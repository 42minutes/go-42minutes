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

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

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

	// TODO(geoah) Handle multiple matched episodes?

	// make sure they are not already in the user's library
	// updates to episodes will be handled differently
	sh, err := d.ulibrary.GetShow(epis[0].ShowID)
	if err != nil && err != minutes.ErrNotFound {
		log.Error("An error occured trying to get show from user library", err)
		return
	}

	if sh != nil {
		log.Error("Show already exists in user library")
		return
	}

	// add show
	show, _ := d.glibrary.GetShow(epis[0].ShowID)
	if err := d.ulibrary.UpsertShow(show); err != nil {
		log.Error("Could not persist show in ulib", err)
	}

	// add seasons
	seas, _ := d.glibrary.GetSeasonsByShow(epis[0].ShowID)
	for _, se := range seas {
		if err := d.ulibrary.UpsertSeason(se); err != nil {
			log.Error("Could not persist season in ulib", err)
		}
		// add episodes
		sepis, err := d.glibrary.GetEpisodesBySeasonNumber(epis[0].ShowID, se.Number)
		if err != nil {
			log.Error("Could not persist episode in ulib", err)
		}
		for _, sep := range sepis {
			d.ulibrary.UpsertEpisode(sep)
		}
	}
}

// Diff will attempt to figure out which episodes are missing from
// the user's library, find their torrents and download them
func (d *daemon) Diff() {
	shows, _ := d.ulibrary.GetShows()
	for _, ushow := range shows {
		gshow, _ := d.glibrary.GetShow(ushow.ID)
		epis, _ := d.differ.Diff(ushow, gshow)
		for _, epi := range epis {
			torr, _ := d.finder.Find(gshow, epi)
			d.downloader.Download(torr[0])
		}
	}
}

func main() {
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
	log.Info(cfg.Rethink.Databases.Library)
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
	diff := &minutes.SimpleDiff{}

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

	// TODO run every x minutes check for missing episodes
	go daem.Diff()

	// start watching for changes
	wtch.Watch(cfg.SeriesPath)
}
