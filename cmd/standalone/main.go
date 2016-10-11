package main

import (
	minutes "github.com/42minutes/go-42minutes"
	trakt "github.com/42minutes/go-trakt"
	"github.com/dancannon/gorethink"
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
}

// HandleWatcherNotification handles watcher notifications
func (d *daemon) HandleWatcherNotification(notifType minutes.NotificationType, file string) {
	// find episode, season, and show
	epis, _ := d.glibrary.QueryEpisodesByFile(file)
	seas, _ := d.glibrary.GetSeason(epis[0].SeasonID)
	show, _ := d.glibrary.GetShow(epis[0].ShowID)
	// make sure they are in the user's library
	d.ulibrary.UpsertShow(show)
	d.ulibrary.UpsertSeason(seas)
	d.ulibrary.UpsertEpisode(epis[0])
}

// Diff will attempt to figure out which episodes are missing from
// the user's library, find their torrents and download them
func (d *daemon) Diff() {
	shows, _ := d.ulibrary.GetShows()
	for _, ushow := range shows {
		gshow, _ := d.glibrary.GetShow(ushow.ID)
		epis, _ := d.differ.Diff(ushow, gshow)
		for _, epi := range epis {
			torr, _ := d.finder.Find(epi)
			d.downloader.Download(torr[0])
		}
	}
}

func main() {
	// trakt.tv client
	trkt := trakt.NewClient(
		"CLIENT_ID",
		trakt.TokenAuth{AccessToken: "ACCESS_TOKEN"},
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

	// standalone daemon
	daem := &daemon{
		glibrary:   glib,
		ulibrary:   ulib,
		finder:     fndr,
		downloader: dwnl,
		differ:     diff,
	}

	// create a new file watcher
	wtch := &minutes.FileWatcher{}
	// notify daemon when something changes
	wtch.Notify(daem)

	// start watching for changes
	go wtch.Watch("/tmp/tvseries", true)

	// TODO run every x minutes check for missing episodes
	go daem.Diff()
}
