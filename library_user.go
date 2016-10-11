package minutes

import "github.com/dancannon/gorethink"

// UserLibrary is a read-write user-specific library
type UserLibrary struct {
	rethinkdb *gorethink.Session
	userID    string
}

// NewUserLibrary returns a UserLibrary
func NewUserLibrary(rethinkdb *gorethink.Session, userID string) *UserLibrary {
	return &UserLibrary{
		rethinkdb: rethinkdb,
		userID:    userID,
	}
}

// UpsertShow adds or updates a show
// or error with ErrNotImplemented, or ErrInternalServer
func (l *UserLibrary) UpsertShow(show *Show) error {
	return nil
}

// UpsertSeason adds or updates a season
// or errors with ErrNotImplemented, or ErrInternalServer, or ErrMissingShow
func (l *UserLibrary) UpsertSeason(season *Season) error {
	return nil
}

// UpsertEpisode adds or updates a episode
// or errors with ErrNotImplemented, or ErrInternalServer, ErrMissingShow
// or ErrMissingSeason
func (l *UserLibrary) UpsertEpisode(episode *Episode) error {
	return nil
}

// GetShow returns a Show
// or errors with ErrNotFound, or ErrInternalServer
func (l *UserLibrary) GetShow(showID string) (*Show, error) {
	return nil, nil
}

// GetShows returns all Shows
// or errors with ErrInternalServer
func (l *UserLibrary) GetShows() ([]*Show, error) {
	return []*Show{}, nil
}

// GetSeason returns a Season
// or errors with ErrNotFound, or ErrInternalServer
func (l *UserLibrary) GetSeason(seasonID string) (*Season, error) {
	return nil, nil
}

// GetSeasonByNumber returns a Season given a Show's ID and a Season number
// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
func (l *UserLibrary) GetSeasonByNumber(showID string, seasonNumber int) (*Season, error) {
	return nil, nil
}

// GetEpisode returns an Episode
// or errors with ErrNotFound, or ErrInternalServer
func (l *UserLibrary) GetEpisode(episodeID string) (*Episode, error) {
	return nil, nil
}

// GetEpisodeByNumber returns an Episode  given a Show's ID a Season number
// and Episode's number
// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
func (l *UserLibrary) GetEpisodeByNumber(showID string, seasonNumber, episodeNumber int) (*Episode, error) {
	return nil, nil
}

// QueryShowsByTitle returns all Shows that match a partial title ordered
// by their probability
// or errors with ErrInternalServer
func (l *UserLibrary) QueryShowsByTitle(title string) ([]*Show, error) {
	return []*Show{}, nil
}

// QueryEpisodesByFile returns all episodes that match a filename or full
// path, ordered by their probability
// or errors with ErrInternalServer
func (l *UserLibrary) QueryEpisodesByFile(filename string) ([]*Episode, error) {
	return []*Episode{}, nil
}
