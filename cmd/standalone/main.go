package main

import (
	minutes "github.com/42minutes/go-42minutes"
	trakt "github.com/42minutes/go-trakt"
	gin "github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("standalone")

// daemon implements the WatchNotifier interface so it can be notified for
// updates
type daemon struct {
	config     *Config
	watcher    minutes.Watcher
	glibrary   minutes.ShowLibrary
	ulibrary   minutes.UserLibrary
	finder     minutes.Finder
	downloader minutes.Downloader
	differ     minutes.Differ
	matcher    minutes.Matcher
	queue      *minutes.Queue
}

// HandleWatcherNotification handles watcher notifications
func (d *daemon) HandleWatcherNotification(notifType minutes.NotificationType, path string) {
	// find episode, season, and show
	epis, _ := d.matcher.Match(path)
	// TODO(geoah) Handle error
	if len(epis) == 0 {
		return
	}

	ep := epis[0]

	// TODO(geoah) Handle multiple matched episodes?

	// make sure they are not already in the user's library
	// updates to episodes will be handled differently
	sh, err := d.ulibrary.GetShow(ep.ShowID)
	// else add or update show
	if err == minutes.ErrNotFound {
		gsh, err := d.glibrary.GetShow(ep.ShowID)
		if err != nil {
			log.Error("Could not get show from glib", err)
			return
		}
		sh = &minutes.UserShow{
			ID:    gsh.ID,
			Title: gsh.Title,
		}
		if err := d.ulibrary.UpsertShow(sh); err != nil {
			log.Error("Could not persist show in ulib", err)
		}
	} else if err != nil && err != minutes.ErrNotFound {
		log.Error("An error occured trying to get show from user library", err)
		return
	}

	se, err := d.ulibrary.GetSeason(sh.ID, ep.Season)
	// else add or update season
	if err == minutes.ErrNotFound {
		gses, _ := d.glibrary.GetSeasons(sh.ID)
		for _, gse := range gses {
			if gse.Number == ep.Season {
				se = &minutes.UserSeason{
					ShowID: sh.ID,
					Number: gse.Number,
				}
				break
			}
		}
		if err := d.ulibrary.UpsertSeason(se); err != nil {
			log.Error("Could not persist season in ulib", err)
		}
	} else if err != nil && err != minutes.ErrNotFound {
		log.Error("An error occured trying to get season from user library", err)
		return
	}

	// else add or update episode
	uep, err := d.ulibrary.GetEpisode(sh.ID, ep.Season, ep.Number)
	if err != nil && err != minutes.ErrNotFound {
		log.Error("An error occured trying to get season from user library", err)
		return
	} else if err == minutes.ErrNotFound {
		gep, _ := d.glibrary.GetEpisode(sh.ID, ep.Season, ep.Number)
		uep = &minutes.UserEpisode{
			ShowID: sh.ID,
			Season: se.Number,
			Number: ep.Number,
			Files:  []*minutes.UserFile{},
		}
		uep.MergeInPlace(gep)
	}

	for _, ufile := range uep.Files {
		// TODO(geoah) Check CRC32
		if ufile.Name == path {
			log.Info("Episode already exists in user's library")
			return
		}
	}

	file := ep.Files[0]
	file.Status = "ok"
	uep.Files = append(uep.Files, file)

	if err := d.ulibrary.UpsertEpisode(uep); err != nil {
		log.Error("Could not persist episode in ulib", err)
		return
	}
	log.Infof(">> Added '%s' S%02dE%02d to user's library", sh.Title, ep.Season, ep.Number)
}

// Diff will attempt to figure out which episodes are missing from
// the user's library, find their torrents and download them
func (d *daemon) Diff() {
	// Add episodes to queue
	log.Info("Running diff")
	shows, _ := d.ulibrary.GetShows()
	for _, ush := range shows {
		log.Infof("> Diffing %s", ush.Title)
		gsh, _ := d.glibrary.GetShow(ush.ID)
		eps, _ := d.differ.Diff(ush, gsh)
		for _, ep := range eps {
			log.Infof(">> Marking %s S%02dE%02d for download", gsh.Title, ep.Season, ep.Number)
			d.queue.Add(ep, &minutes.UserFile{
				Status: "pending",
			})
		}
	}
}

func main() {
	logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	log.Info("Reading config file.")
	cfg, err := loadConfig("./config.json")
	if err != nil {
		log.Fatal("Could not load config file.", err)
	}

	// trakt.tv client
	trkt := trakt.NewClient(
		cfg.Trakt.ClientID,
		trakt.TokenAuth{AccessToken: ""},
	)

	// global ro trakt library
	glib := minutes.NewTraktLibrary(trkt)

	// SQLite instance for UserLibrary
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal("Could not init SQLite")
	}

	// user rw library for single hardcoded user id
	ulib, err := minutes.NewSqlUserLibrary(db)
	if err != nil {
		log.Fatal("Could not init UserLibrary", err)
	}
	defer ulib.Close()

	// torrent finder
	fndr := &minutes.TorrentFinder{}

	// torrent download manager
	dwnl := minutes.NewTorrentDownloader(cfg.WatchPath)

	// simple differ
	diff := minutes.NewSimpleDiff(ulib, glib)

	// simple matcher
	mtch, _ := minutes.NewSimpleMatch(glib)

	// create a new file watcher
	wtch := &minutes.FileWatcher{}

	// queue
	qu, _ := minutes.NewQueue(db, fndr, glib, ulib, dwnl)

	// standalone daemon
	daem := &daemon{
		config:     cfg,
		glibrary:   glib,
		ulibrary:   ulib,
		finder:     fndr,
		downloader: dwnl,
		differ:     diff,
		matcher:    mtch,
		watcher:    wtch,
		queue:      qu,
	}

	// notify daemon when something changes
	wtch.Notify(daem)

	// start processing the queue
	qu.Process()

	// create api for http
	api := minutes.NewAPI(glib, ulib)

	// start http server
	app := gin.Default()
	app.GET("/shows", api.HandleShows)
	app.POST("/shows", api.HandleShowPost)
	app.GET("/shows/:show_id", api.HandleShow)
	app.GET("/shows/:show_id/seasons", api.HandleSeasons)
	app.GET("/shows/:show_id/seasons/:season", api.HandleSeason)
	app.GET("/shows/:show_id/seasons/:season/episodes", api.HandleEpisodes)
	app.GET("/shows/:show_id/seasons/:season/episodes/:episode", api.HandleEpisode)
	go app.Run(":8081") // TODO(geoah) Make port configurable

	// start shell
	daem.startShell()
}
