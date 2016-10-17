package minutes

import rethink "github.com/dancannon/gorethink"

const (
	tableShows    = "shows"
	tableSeasons  = "seasons"
	tableEpisodes = "episodes"
)

// UserLibrary is a read-write user-specific library
type UserLibrary struct {
	rethinkdb *rethink.Session
	userID    string
}

// NewUserLibrary returns a UserLibrary
func NewUserLibrary(redb *rethink.Session, uid string) *UserLibrary {
	return &UserLibrary{
		rethinkdb: redb,
		userID:    uid,
	}
}

// UpsertShow adds or updates a show
// or error with ErrNotImplemented, or ErrInternalServer
func (l *UserLibrary) UpsertShow(show *Show) error {
	return l.upsert(tableShows, show)
}

// UpsertSeason adds or updates a season
// or errors with ErrNotImplemented, or ErrInternalServer, or ErrMissingShow
func (l *UserLibrary) UpsertSeason(season *Season) error {
	return l.upsert(tableSeasons, season)
}

// UpsertEpisode adds or updates a episode
// or errors with ErrNotImplemented, or ErrInternalServer, ErrMissingShow
// or ErrMissingSeason
func (l *UserLibrary) UpsertEpisode(episode *Episode) error {
	return l.upsert(tableEpisodes, episode)
}

// GetShow returns a Show
// or errors with ErrNotFound, or ErrInternalServer
func (l *UserLibrary) GetShow(id string) (*Show, error) {
	sh := &Show{}
	if err := l.get(tableShows, id, sh); err != nil {
		return nil, err
	}
	return sh, nil
}

// GetShows returns all Shows
// or errors with ErrInternalServer
func (l *UserLibrary) GetShows() ([]*Show, error) {
	res, err := rethink.Table(tableShows).Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	shs := []*Show{}
	if res.IsNil() {
		return shs, nil
	}
	err = res.All(&shs)
	if err != nil {
		return nil, ErrInternalServer
	}
	return shs, nil
}

// GetSeason returns a Season
// or errors with ErrNotFound, or ErrInternalServer
func (l *UserLibrary) GetSeason(id string) (*Season, error) {
	se := &Season{}
	err := l.get(tableSeasons, id, se)
	return se, err
}

// GetSeasonsByShow returns all Seasons for a show
// or errors with ErrNotFound, or ErrInternalServer
func (l *UserLibrary) GetSeasonsByShow(sid string) ([]*Season, error) {
	qr := rethink.Table(tableSeasons).GetAllByIndex("show_id", sid)
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	ses := []*Season{}
	if res.IsNil() {
		return ses, nil
	}
	err = res.All(&ses)
	if err != nil {
		return nil, ErrInternalServer
	}
	return ses, nil
}

// GetSeasonByNumber returns a Season given a Show's ID and a Season number
// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
func (l *UserLibrary) GetSeasonByNumber(sid string, sn int) (*Season, error) {
	qr := rethink.Table(tableSeasons).GetAllByIndex("show_id", sid)
	qr = qr.Filter(map[string]interface{}{"number": sn})
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	if res.IsNil() {
		return nil, ErrNotFound
	}
	se := &Season{}
	if err := res.One(se); err != nil {
		return nil, ErrInternalServer
	}
	return se, nil
}

// GetEpisode returns an Episode
// or errors with ErrNotFound, or ErrInternalServer
func (l *UserLibrary) GetEpisode(id string) (*Episode, error) {
	ep := &Episode{}
	err := l.get(tableEpisodes, id, ep)
	return ep, err
}

// GetEpisodesBySeason returns all Episodes for a Season
// or errors with ErrNotFound, or ErrInternalServer
func (l *UserLibrary) GetEpisodesBySeason(sid string) ([]*Episode, error) {
	qr := rethink.Table(tableEpisodes).GetAllByIndex("season_id", sid)
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	eps := []*Episode{}
	if res.IsNil() {
		return eps, nil
	}
	err = res.All(&eps)
	if err != nil {
		return nil, ErrInternalServer
	}
	return eps, nil
}

// GetEpisodesBySeasonNumber returns all Shows for a show and season number
// or errors with ErrNotFound, or ErrInternalServer
func (l *UserLibrary) GetEpisodesBySeasonNumber(sid string, sn int) ([]*Episode, error) {
	qr := rethink.Table(tableEpisodes).GetAllByIndex("show_id", sid)
	qr = qr.Filter(map[string]interface{}{"season": sn})
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	eps := []*Episode{}
	if res.IsNil() {
		return eps, nil
	}
	err = res.All(&eps)
	if err != nil {
		return nil, ErrInternalServer
	}
	return eps, nil
}

// GetEpisodeByNumber returns an Episode  given a Show's ID a Season number
// and Episode's number
// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
func (l *UserLibrary) GetEpisodeByNumber(sid string, sn, en int) (*Episode, error) {
	qr := rethink.Table(tableEpisodes).GetAllByIndex("show_id", sid)
	qr = qr.Filter(map[string]interface{}{"season": sn, "number": en})
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		log.Info(err)
		return nil, ErrInternalServer
	}
	defer res.Close()
	if res.IsNil() {
		return nil, ErrNotFound
	}
	ep := &Episode{}
	if err := res.One(ep); err != nil {
		return nil, ErrInternalServer
	}
	return ep, nil
}

// QueryShowsByTitle returns all Shows that match a partial title ordered
// by their probability
// or errors with ErrInternalServer
func (l *UserLibrary) QueryShowsByTitle(title string) ([]*Show, error) {
	qr := rethink.Table(tableShows)
	qr = qr.Filter(rethink.Row.Field("title").Match(title))
	res, err := qr.Run(l.rethinkdb)
	if err != nil {
		return nil, ErrInternalServer
	}
	defer res.Close()
	shs := []*Show{}
	if res.IsNil() {
		return shs, nil
	}
	err = res.All(&shs)
	if err != nil {
		return nil, ErrInternalServer
	}
	return shs, nil
}

func (l *UserLibrary) upsert(tbl string, doc interface{}) error {
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

func (l *UserLibrary) get(tbl, id string, doc interface{}) error {
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
