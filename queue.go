package minutes

import (
	"time"

	torrentlookup "github.com/42minutes/go-torrentlookup"
	rethink "github.com/dancannon/gorethink"
)

var retryAfterHours = 3

// Queue -
type Queue struct {
	rethinkdb *rethink.Session
	finder    Finder
	glibrary  ShowLibrary
	ulibrary  UserLibrary
	dwnl      Downloader
}

// NewQueue -
func NewQueue(redb *rethink.Session, fndr Finder, glib ShowLibrary, ulib UserLibrary, dwnl Downloader) (*Queue, error) {
	return &Queue{
		rethinkdb: redb,
		finder:    fndr,
		glibrary:  glib,
		ulibrary:  ulib,
		dwnl:      dwnl,
	}, nil
}

// Add -
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

// Process -
func (q *Queue) Process() {
	go func() {
		for {
			items, err := q.ulibrary.QueryEpisodesForFinder()
			if err != nil {
				log.Error("Could not get items", err)
			}

			log.Warningf("[process-lookup] Processing %d episodes.", len(items))

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
					if err != ErrNotFound {
						log.Errorf("[process-lookup] Error trying to find magnet for %s, %v", it, err)
					}
				} else {
					if len(down) > 0 {
						it.Infohash = down[0].GetID()
						log.Infof("[process-lookup] Found hash for magnet %s", it)
						if err := q.ulibrary.UpsertEpisode(it); err != nil {
							log.Error("[process-lookup] Could not update episode.", err)
						}
					}
				}
				if it.Infohash == "" {
					log.Infof("[process-lookup] Could not find magnet for %s, will retry in %d hours", it, retryAfterHours)
					it.RetryDatetime = time.Now().Add(time.Hour * time.Duration(retryAfterHours)).UTC().Unix()
					if err := q.ulibrary.UpsertEpisode(it); err != nil {
						log.Error("[process-lookup] Could not update episode with retry.", err)
					}
				}
			}

			time.Sleep(time.Second * 30)
		}
	}()
	go func() {
		for {
			items, err := q.ulibrary.QueryEpisodesForDownloader()
			if err != nil {
				log.Error("[process-download] Could not get items", err)
			}

			log.Warningf("[process-download] Processing %d episodes.", len(items))

			for _, it := range items {
				err := q.dwnl.Download(&MagnetDownloadable{
					Infohash: it.Infohash,
					Magnet:   torrentlookup.CreateFakeMagnet(it.Infohash),
				})

				if err != nil {
					if err == ErrNotFound {
						log.Errorf("[process-download] Error trying to download %s %v", it, err)
					}
				} else {
					it.Downloaded = true
					if err := q.ulibrary.UpsertEpisode(it); err != nil {
						log.Error("[process-download] Could not update episode after downloading.", err)
					}
					log.Infof("[process-download] Downloaded %s ", it)
				}
			}
			time.Sleep(time.Second * 45)
		}
	}()
}
