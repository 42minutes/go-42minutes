package minutes

import (
	"time"

	torrentlookup "github.com/42minutes/go-torrentlookup"
	rethink "github.com/dancannon/gorethink"
)

type Queue struct {
	rethinkdb *rethink.Session
	finder    Finder
	glibrary  ShowLibrary
	ulibrary  UserLibrary
	dwnl      Downloader
}

func NewQueue(redb *rethink.Session, fndr Finder, glib ShowLibrary, ulib UserLibrary, dwnl Downloader) (*Queue, error) {
	return &Queue{
		rethinkdb: redb,
		finder:    fndr,
		glibrary:  glib,
		ulibrary:  ulib,
		dwnl:      dwnl,
	}, nil
}

func (q *Queue) Add(ep *Episode) error {
	uep := &UserEpisode{
		ShowID:        ep.ShowID,
		Season:        ep.Season,
		Number:        ep.Number,
		Downloaded:    false,
		RetryDatetime: time.Now().Add(-10 * time.Minute).UTC().Unix(),
	}
	if err := q.ulibrary.UpsertEpisode(uep); err != nil {
		return ErrInternalServer
	}
	return nil
}

func (q *Queue) Process() {
	go func() {
		for {
			items, err := q.ulibrary.QueryEpisodesForFinder()
			if err != nil {
				log.Error("Could not get items", err)
			}

			log.Warningf(">>> Processing (lookup) %d episodes.", len(items))

			for _, it := range items {
				ep, err := q.glibrary.GetEpisode(it.ShowID, it.Season, it.Number)
				if err != nil {
					log.Error(err)
				}

				sh, err := q.glibrary.GetShow(it.ShowID)
				if err != nil {
					log.Error(err)
				}

				down, err := q.finder.Find(sh, ep)
				if err != nil {
					log.Warning(">>> Could not find magnet", err)
				} else {
					if len(down) > 0 {
						log.Infof(">>> Found hash for magnet: %s", down[0].GetID())
						it.Infohash = down[0].GetID()
						if err := q.ulibrary.UpsertEpisode(it); err != nil {
							// TODO(geoah) log error
							log.Error(err)
						}
					}
				}
			}

			time.Sleep(time.Second * 30)
		}
	}()
	go func() {
		for {
			log.Warning(">>> Downloading queue...")
			items, err := q.ulibrary.QueryEpisodesForDownloader()
			if err != nil {
				log.Error("Could not get items", err)
			}

			log.Warningf(">>> Processing (downloading) %d episodes.", len(items))

			for _, it := range items {
				err := q.dwnl.Download(&MagnetDownloadable{
					Infohash: it.Infohash,
					Magnet:   torrentlookup.CreateFakeMagnet(it.Infohash),
				})

				if err != nil {
					log.Error("Could not download", err)
				} else {
					it.Downloaded = true
					if err := q.ulibrary.UpsertEpisode(it); err != nil {
						// TODO(geoah) log error
						log.Error(err)
					}
					log.Infof(">>> Downloaded ShowID: %s Season: %d Episode: %d Infohash: %s",
						it.ShowID, it.Season, it.Number, it.Infohash)
				}
			}
			time.Sleep(time.Second * 45)
		}
	}()
}
