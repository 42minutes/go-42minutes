package minutes

import (
	"fmt"
	"time"

	rethink "github.com/dancannon/gorethink"
)

const (
	tableQueue = "queue"
)

type item struct {
	episode       string `gorethink:"episode_id"`
	infohash      string `gorethink:"infohash"`
	downloaded    bool   `gorethink:"downloaded"`
	retryDatetime int64  `gorethink:"retry_time"`
}

type Queue struct {
	rethinkdb *rethink.Session
	finder    Finder
}

func NewQueue(redb *rethink.Session, fndr Finder) (*Queue, error) {
	return &Queue{
		rethinkdb: redb,
		finder:    fndr,
	}, nil
}

func (q *Queue) Add(ep *Episode) error {
	// Add episode to queue
	insertOpts := rethink.InsertOpts{
		Conflict: "update",
	}

	it := item{
		episode:       ep.ID,
		infohash:      "",
		downloaded:    false,
		retryDatetime: time.Now().Add(-10 * time.Minute).UTC().Unix(),
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
			res, err := rethink.Table(tableQueue).Filter(rethink.Row.Field("retry_time").Le(time.Now().UTC().Unix())).Run(q.rethinkdb)
			if err != nil {
				log.Error(err)
			}

			fmt.Println(res)
			// find infohashes for episodes that retry is in the past
			// update episode with infohash

			time.Sleep(time.Minute * 5)
		}
	}()
	go func() {
		for {
			// get all episodes that have infohashes
			// downloaded = false
			// download
			// update download
		}
		time.Sleep(time.Minute * 2)
	}()
}
