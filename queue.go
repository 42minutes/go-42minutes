package minutes

import (
	"time"

	torrentlookup "github.com/42minutes/go-torrentlookup"
	"github.com/jinzhu/gorm"
)

var retryAfterHours = 3

// Queue -
type Queue struct {
	sqldb    *gorm.DB
	finder   Finder
	glibrary ShowLibrary
	ulibrary UserLibrary
	dwnl     Downloader
}

// NewQueue -
func NewQueue(db *gorm.DB, fndr Finder, glib ShowLibrary, ulib UserLibrary, dwnl Downloader) (*Queue, error) {
	return &Queue{
		sqldb:    db,
		finder:   fndr,
		glibrary: glib,
		ulibrary: ulib,
		dwnl:     dwnl,
	}, nil
}

// Add -
func (q *Queue) Add(ep *Episode, file *UserFile) error {
	// TODO(geoah) Check if we need to start locking
	uep, err := q.ulibrary.GetEpisode(ep.ShowID, ep.Season, ep.Number)
	if err != nil && err != ErrNotFound {
		return ErrInternalServer
	} else if err == ErrNotFound {
		uep = &UserEpisode{
			ShowID: ep.ShowID,
			Season: ep.Season,
			Number: ep.Number,
			Files:  []*UserFile{},
		}
	}
	for _, ufile := range uep.Files {
		if file.Resolution == ufile.Resolution && file.Source == ufile.Source {
			return nil
		}
	}
	uep.Files = append(uep.Files, file)
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
					log.Error("%v", err)
				}

				sh, err := q.glibrary.GetShow(it.ShowID)
				if err != nil {
					log.Error("%v", err)
				}

				down, err := q.finder.Find(sh, ep)
				if err != nil {
					if err != ErrNotFound {
						log.Errorf("[process-lookup] Error trying to find magnet for %s, %v", it, err)
					}
				} else {
					if len(down) > 0 {
						it.Files[0].Status = "found"
						it.Files[0].Infohash = down[0].GetID()
						log.Infof("[process-lookup] Found hash for magnet %s", it)
						if err := q.ulibrary.UpsertEpisode(it); err != nil {
							log.Error("[process-lookup] Could not update episode.", err)
						}
					}
				}
				if it.Files[0].Infohash == "" {
					log.Infof("[process-lookup] Could not find magnet for %s, will retry in %d hours", it, retryAfterHours)
					it.Files[0].RetryTime = time.Now().Add(time.Hour * time.Duration(retryAfterHours)).UTC().Unix()
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
					Infohash: it.Files[0].Infohash,
					Magnet:   torrentlookup.CreateFakeMagnet(it.Files[0].Infohash),
				})

				if err != nil {
					if err == ErrNotFound {
						log.Errorf("[process-download] Error trying to download %s %v", it, err)
					}
				} else {
					it.Files[0].Status = "downloading"
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
