package minutes

import (
	"time"

	rethink "github.com/dancannon/gorethink"
)

const (
	tableQueue = "queue"
)

type item struct {
	episode       *Episode  `gorethink:"episode"`
	infohash      string    `gorethink:"infohash"`
	downloaded    bool      `gorethink:"downloaded"`
	retryDatetime time.Time `gorethink:"retryDatetime"`
}

type Queue struct {
	rethinkdb *rethink.Session
	finder    Finder
}

func NewQueue(redb *rethink.Session, fndr Finder) *Queue {
	return &Queue{
		rethinkdb: redb,
		finder:    fndr,
	}
}

func (q *Queue) Add(ep *Episode) error {
	// Add episode to queue
	insertOpts := rethink.InsertOpts{
		Conflict: "update",
	}

	it := item{
		episode:       ep,
		infohash:      "",
		downloaded:    false,
		retryDatetime: time.Now().Add(-10 * time.Minute),
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
			// run dirr
		}
		time.Sleep(time.Minute * 5)
	}()
	go func() {
		for {
			// find all episodes for which we need to find torrents
		}
		time.Sleep(time.Minute * 2)
	}()
}
