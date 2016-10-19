package minutes

import (
	"time"

	torrentlookup "github.com/42minutes/go-torrentlookup"
	rethink "github.com/dancannon/gorethink"
)

const (
	tableQueue = "queue"
)

type item struct {
	ShowID        string `gorethink:"show_id"`
	SeasonNumber  int    `gorethink:"season_number"`
	EpisodeNumber int    `gorethink:"episode_number"`
	Infohash      string `gorethink:"infohash"`
	Downloaded    bool   `gorethink:"downloaded"`
	RetryDatetime int64  `gorethink:"retry_time"`
}

type Queue struct {
	rethinkdb *rethink.Session
	finder    Finder
	glib      Library
	dwnl      Downloader
}

func NewQueue(redb *rethink.Session, fndr Finder, glib Library, dwnl Downloader) (*Queue, error) {
	return &Queue{
		rethinkdb: redb,
		finder:    fndr,
		glib:      glib,
		dwnl:      dwnl,
	}, nil
}

func (q *Queue) Add(ep *Episode) error {
	// Add episode to queue
	insertOpts := rethink.InsertOpts{
		Conflict: "update",
	}

	it := item{
		ShowID:        ep.ShowID,
		SeasonNumber:  ep.Season,
		EpisodeNumber: ep.Number,
		Infohash:      "",
		Downloaded:    false,
		RetryDatetime: time.Now().Add(-10 * time.Minute).UTC().Unix(),
	}

	qr := rethink.Table(tableQueue).Insert(it, insertOpts)
	if _, err := qr.RunWrite(q.rethinkdb); err != nil {
		return ErrInternalServer
	}
	return nil
}

func (q *Queue) Process() {
	go func() {
		for {
			// get episodes from queue from db
			res, err := rethink.Table(tableQueue).Filter(
				rethink.And(
					rethink.Row.Field("retry_time").Le(time.Now().UTC().Unix()),
					rethink.Row.Field("infohash").Eq(""),
				),
			).Run(q.rethinkdb)
			if err != nil {
				log.Error(err)
			}

			var items = []item{}
			if err := res.All(&items); err != nil {
				log.Error("Could not get items", err)
			}

			for _, it := range items {
				ep, err := q.glib.GetEpisodeByNumber(it.ShowID, it.SeasonNumber, it.EpisodeNumber)
				if err != nil {
					log.Error(err)
				}

				sh, err := q.glib.GetShow(it.ShowID)
				if err != nil {
					log.Error(err)
				}

				down, err := q.finder.Find(sh, ep)
				if err != nil {
					log.Warning(">>> Could not find magnet", err)
				} else {
					if len(down) > 0 {
						log.Infof(">>> Found hash for magnet: %s", down[0].GetID())
						insertOpts := rethink.InsertOpts{
							Conflict: "update",
						}
						it.Infohash = down[0].GetID()

						qr := rethink.Table(tableQueue).Insert(it, insertOpts)
						if _, err := qr.RunWrite(q.rethinkdb); err != nil {
							// TODO(geoah) log error
							log.Error(err)
						}
					}
				}
			}

			time.Sleep(time.Minute * 2)
		}
	}()
	go func() {
		for {
			res, err := rethink.Table(tableQueue).Filter(
				rethink.And(
					rethink.Row.Field("infohash").Ne(""),
					rethink.Row.Field("downloaded").Eq(false),
				),
			).Run(q.rethinkdb)
			if err != nil {
				log.Error(err)
			}

			var items = []item{}
			if err := res.All(&items); err != nil {
				log.Error("Could not get items", err)
			}

			for _, it := range items {
				err := q.dwnl.Download(&MagnetDownloadable{
					Infohash: it.Infohash,
					Magnet:   torrentlookup.CreateFakeMagnet(it.Infohash),
				})

				if err != nil {
					log.Error("Could not download", err)
				} else {
					it.Downloaded = true
					insertOpts := rethink.InsertOpts{
						Conflict: "update",
					}

					qr := rethink.Table(tableQueue).Insert(it, insertOpts)
					if _, err := qr.RunWrite(q.rethinkdb); err != nil {
						// TODO(geoah) log error
						log.Error(err)
					}
					log.Infof(">>> Downloaded ShowID: %s Season: %d Episode: %d Infohash: %s", it.ShowID, it.SeasonNumber, it.EpisodeNumber, it.Infohash)
				}
			}
			time.Sleep(time.Minute * 2)
		}
	}()
}
