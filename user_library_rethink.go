package minutes

import (
	"time"

	rethink "github.com/dancannon/gorethink"
)

const (
	tableShows    = "shows"
	tableSeasons  = "seasons"
	tableEpisodes = "episodes"
)

// RethinkUserLibrary is a read-write user-specific library
type RethinkUserLibrary struct {
	rethinkdb *rethink.Session
	userID    string
}

// NewUserLibrary returns a RethinkUserLibrary
func NewUserLibrary(redb *rethink.Session, uid string) *RethinkUserLibrary {
	return &RethinkUserLibrary{
		rethinkdb: redb,
		userID:    uid,
	}
}

// UpsertShow adds or updates a show
// or error with ErrNotImplemented, or ErrInternalServer
func (l *RethinkUserLibrary) UpsertShow(show *UserShow) error {
	return l.upsert(tableShows, show)
}

// UpsertSeason adds or updates a season
// or errors with ErrNotImplemented, or ErrInternalServer, or ErrMissingShow
func (l *RethinkUserLibrary) UpsertSeason(season *UserSeason) error {
	season.CID = season.GetCID()
	return l.upsert(tableSeasons, season)
}

// UpsertEpisode adds or updates a episode
// or errors with ErrNotImplemented, or ErrInternalServer, ErrMissingShow
// or ErrMissingSeason
func (l *RethinkUserLibrary) UpsertEpisode(episode *UserEpisode) error {
	episode.CID = episode.GetCID()
	return l.upsert(tableEpisodes, episode)
}

// GetShow returns a UserShow
// or errors with ErrNotFound, or ErrInternalServer
func (l *RethinkUserLibrary) GetShow(id string) (*UserShow, error) {
	sh := &UserShow{}
	if err := l.get(tableShows, id, sh); err != nil {
		return nil, err
	}
	return sh, nil
}

// GetShows returns all Shows
// or errors with ErrInternalServer
func (l *RethinkUserLibrary) GetShows() ([]*UserShow, error) {
	res, err := rethink.Table(tableShows).Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	shs := []*UserShow{}
	if res.IsNil() {
		return shs, nil
	}
	err = res.All(&shs)
	if err != nil {
		return nil, ErrInternalServer
	}
	return shs, nil
}

// GetSeasons returns all Seasons for a show
// or errors with ErrNotFound, or ErrInternalServer
func (l *RethinkUserLibrary) GetSeasons(sid string) ([]*UserSeason, error) {
	qr := rethink.Table(tableSeasons)
	qr = qr.Filter(map[string]interface{}{"show_id": sid})
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	ses := []*UserSeason{}
	if res.IsNil() {
		return ses, nil
	}
	err = res.All(&ses)
	if err != nil {
		return nil, ErrInternalServer
	}
	return ses, nil
}

// GetSeason returns a UserSeason given a UserShow's ID and a UserSeason number
// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
func (l *RethinkUserLibrary) GetSeason(sid string, sn int) (*UserSeason, error) {
	qr := rethink.Table(tableSeasons).Get([]interface{}{sid, sn})
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	if res.IsNil() {
		return nil, ErrNotFound
	}
	se := &UserSeason{}
	if err := res.One(se); err != nil {
		return nil, ErrInternalServer
	}
	return se, nil
}

// GetEpisodes returns all Shows for a show and season number
// or errors with ErrNotFound, or ErrInternalServer
func (l *RethinkUserLibrary) GetEpisodes(sid string, sn int) ([]*UserEpisode, error) {
	qr := rethink.Table(tableEpisodes)
	qr = qr.Filter(map[string]interface{}{"show_id": sid, "season": sn})
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	eps := []*UserEpisode{}
	if res.IsNil() {
		return eps, nil
	}
	err = res.All(&eps)
	if err != nil {
		return nil, ErrInternalServer
	}
	return eps, nil
}

// GetEpisode returns an UserEpisode given a UserShow's ID a UserSeason number
// and UserEpisode's number
// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
func (l *RethinkUserLibrary) GetEpisode(sid string, sn, en int) (*UserEpisode, error) {
	qr := rethink.Table(tableEpisodes).Get([]interface{}{sid, sn, en})
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		log.Info(err)
		return nil, ErrInternalServer
	}
	defer res.Close()
	if res.IsNil() {
		return nil, ErrNotFound
	}
	ep := &UserEpisode{}
	if err := res.One(ep); err != nil {
		return nil, ErrInternalServer
	}
	return ep, nil
}

// QueryShowsByTitle returns all Shows that match a partial title ordered
// by their probability
// or errors with ErrInternalServer
func (l *RethinkUserLibrary) QueryShowsByTitle(title string) ([]*UserShow, error) {
	qr := rethink.Table(tableShows)
	qr = qr.Filter(rethink.Row.Field("title").Match(title))
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	shs := []*UserShow{}
	if res.IsNil() {
		return shs, nil
	}
	err = res.All(&shs)
	if err != nil {
		return nil, ErrInternalServer
	}
	return shs, nil
}

func (l *RethinkUserLibrary) upsert(tbl string, doc interface{}) error {
	insertOpts := rethink.InsertOpts{
		Conflict: "update",
	}
	qr := rethink.Table(tbl).Insert(doc, insertOpts)
	if _, err := qr.RunWrite(l.rethinkdb); err != nil {
		// TODO(geoah) log error
		return ErrInternalServer
	}
	return nil
}

func (l *RethinkUserLibrary) get(tbl, id string, doc interface{}) error {
	res, err := rethink.Table(tbl).Get(id).Run(l.rethinkdb)
	if err != nil {
		return ErrInternalServer
	}
	defer res.Close()
	if res.IsNil() {
		return ErrNotFound
	}
	if err := res.One(doc); err != nil {
		return ErrInternalServer
	}
	return nil
}

// QueryEpisodesForFinder -
func (l *RethinkUserLibrary) QueryEpisodesForFinder() ([]*UserEpisode, error) {
	res, err := rethink.Table(tableEpisodes).Filter(
		rethink.And(
			rethink.Row.Field("downloaded").Eq(false),
			rethink.Row.Field("retry_time").Le(time.Now().UTC().Unix()),
			rethink.Row.Field("infohash").Eq(""),
		),
	).Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	eps := []*UserEpisode{}
	if res.IsNil() {
		return eps, nil
	}
	err = res.All(&eps)
	if err != nil {
		return nil, ErrInternalServer
	}
	return eps, nil
}

// QueryEpisodesForDownloader -
func (l *RethinkUserLibrary) QueryEpisodesForDownloader() ([]*UserEpisode, error) {
	res, err := rethink.Table(tableEpisodes).Filter(
		rethink.And(
			rethink.Row.Field("downloaded").Eq(false),
			rethink.Row.Field("retry_time").Le(time.Now().UTC().Unix()),
			rethink.Row.Field("infohash").Ne(""),
		),
	).Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	eps := []*UserEpisode{}
	if res.IsNil() {
		return eps, nil
	}
	err = res.All(&eps)
	if err != nil {
		return nil, ErrInternalServer
	}
	return eps, nil
}
