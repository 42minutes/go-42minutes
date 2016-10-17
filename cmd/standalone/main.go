package main

import (
	"context"
	"os"

	minutes "github.com/42minutes/go-42minutes"
	trakt "github.com/42minutes/go-trakt"
	"github.com/dancannon/gorethink"
	logging "github.com/op/go-logging"
	"golang.org/x/oauth2"
)

var log = logging.MustGetLogger("standalone")

// daemon implements the WatchNotifier interface so it can be notified for
// updates
type daemon struct {
	watcher    minutes.Watcher
	glibrary   minutes.Library
	ulibrary   minutes.Library
	finder     minutes.Finder
	downloader minutes.Downloader
	differ     minutes.Differ
	matcher    minutes.Matcher
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
		sh, _ = d.glibrary.GetShow(ep.ShowID)
		if err := d.ulibrary.UpsertShow(sh); err != nil {
			log.Error("Could not persist show in ulib", err)
		}
	} else if err != nil && err != minutes.ErrNotFound {
		log.Error("An error occured trying to get show from user library", err)
		return
	}

	se, err := d.ulibrary.GetSeasonByNumber(sh.ID, ep.Season)
	// else add or update season
	if err == minutes.ErrNotFound {
		ses, _ := d.glibrary.GetSeasonsByShow(sh.ID)
		for _, ise := range ses {
			if ise.Number == ep.Season {
				se = ise
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

	_, err = d.ulibrary.GetEpisodeByNumber(sh.ID, ep.Season, ep.Number)
	// else add or update season
	if err == minutes.ErrNotFound {
		uep, _ := d.glibrary.GetEpisodeByNumber(sh.ID, ep.Season, ep.Number)
		if err := d.ulibrary.UpsertEpisode(uep); err != nil {
			log.Error("Could not persist episode in ulib", err)
		}
		log.Infof(">> Added '%s' S%02dE%02d to user's library", sh.Title, ep.Season, ep.Number)
	} else if err != nil && err != minutes.ErrNotFound {
		log.Error("An error occured trying to get season from user library", err)
		return
	}

}

// Diff will attempt to figure out which episodes are missing from
// the user's library, find their torrents and download them
func (d *daemon) Diff() {
	log.Info("Running diff")
	shows, _ := d.ulibrary.GetShows()
	for _, ush := range shows {
		log.Infof("> Diffing %s", ush.Title)
		gsh, _ := d.glibrary.GetShow(ush.ID)
		eps, _ := d.differ.Diff(ush, gsh)
		for _, ep := range eps {
			log.Infof(">> Trying to find way to download %s S%02dE%02d", gsh.Title, ep.Season, ep.Number)
			dnls, _ := d.finder.Find(gsh, ep)
			if len(dnls) > 0 {
				log.Infof(">>> Found hash for magnet: %s", dnls[0].GetID())
				// d.downloader.Download(dnls[0])
			} else {
				log.Infof(">>> Could not find magnet")
			}
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
		log.Info("Could not load config file.", err)
		os.Exit(1)
	}

	log.Info("Getting trakt tokens.")
	oac := &oauth2.Config{
		ClientID:     cfg.Trakt.ClientID,
		ClientSecret: cfg.Trakt.ClientSecret,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.trakt.tv/oauth/authorize",
			TokenURL: "https://api.trakt.tv/oauth/token",
		},
	}

	ctx := context.Background()
	tok := newOAuthToken(ctx, oac)

	log.Info("Got trakt access refresh tokens.", tok.AccessToken, tok.RefreshToken)

	// trakt.tv client
	trkt := trakt.NewClient(
		cfg.Trakt.ClientID,
		trakt.TokenAuth{AccessToken: tok.AccessToken},
	)

	// global ro trakt library
	glib := minutes.NewTraktLibrary(trkt)
	// rethinkdb session for user library
	redb, _ := gorethink.Connect(gorethink.ConnectOpts{
		Address:  "localhost",
		Database: cfg.Rethink.Databases.Library,
	})

	// user rw library for single hardcoded user id
	ulib := minutes.NewUserLibrary(redb, "me")

	// torrent finder
	fndr := &minutes.TorrentFinder{}

	// torrent download manager
	dwnl := &minutes.DownloaderTorrent{}

	// simple differ
	diff := minutes.NewSimpleDiff(ulib, glib)

	// simple matcher
	mtch, _ := minutes.NewSimpleMatch(glib)

	// standalone daemon
	daem := &daemon{
		glibrary:   glib,
		ulibrary:   ulib,
		finder:     fndr,
		downloader: dwnl,
		differ:     diff,
		matcher:    mtch,
	}

	// create a new file watcher
	wtch := &minutes.FileWatcher{}
	// notify daemon when something changes
	wtch.Notify(daem)

	// start watching for changes
	wtch.Watch(cfg.SeriesPath)

	// TODO run every x minutes check for missing episodes
	daem.Diff()
}
