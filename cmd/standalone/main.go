package main

import (
	"context"
	"log"
	"os"

	"fmt"

	minutes "github.com/42minutes/go-42minutes"
	trakt "github.com/42minutes/go-trakt"
	"github.com/dancannon/gorethink"
	"golang.org/x/oauth2"
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
	fmt.Println(notifType, path)
	// find episode, season, and show
	epis, _ := d.matcher.Match(path)
	// TODO(geoah) Implement actual flow
	fmt.Println("Matched", path, "to", epis)
	// seas, _ := d.glibrary.GetSeason(epis[0].SeasonID)
	// show, _ := d.glibrary.GetShow(epis[0].ShowID)
	// make sure they are in the user's library
	// d.ulibrary.UpsertShow(show)
	// d.ulibrary.UpsertSeason(seas)
	// d.ulibrary.UpsertEpisode(epis[0])
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
	log.Println("Reading config file.")
	cfg, err := loadConfig("./config.json")
	if err != nil {
		log.Println("Could not load config file.", err)
		os.Exit(1)
	}

	log.Println("Getting trakt tokens.")
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

	log.Println("Got trakt access refresh tokens.", tok.AccessToken, tok.RefreshToken)

	// trakt.tv client
	trkt := trakt.NewClient(
		cfg.Trakt.ClientID,
		trakt.TokenAuth{AccessToken: tok.AccessToken},
	)

	// global ro trakt library
	glib := minutes.NewTraktLibrary(trkt)

	// rethinkdb session
	redb, _ := gorethink.Connect(gorethink.ConnectOpts{
		Address: "localhost",
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
