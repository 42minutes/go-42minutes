package minutes

import trakt "github.com/42minutes/go-trakt"

// TraktLibrary is a read-only Trakt.tv library
type TraktLibrary struct {
	client *trakt.Client
}

// NewTraktLibrary returns a TraktLibrary
func NewTraktLibrary(client *trakt.Client) *TraktLibrary {
	return &TraktLibrary{
		client: client,
	}
}

// UpsertShow returns ErrNotImplemented
func (l *TraktLibrary) UpsertShow(show *Show) error {
	return ErrNotImplemented
}

// UpsertSeason returns ErrNotImplemented
func (l *TraktLibrary) UpsertSeason(season *Season) error {
	return ErrNotImplemented
}

// UpsertEpisode returns ErrNotImplemented
func (l *TraktLibrary) UpsertEpisode(episode *Episode) error {
	return ErrNotImplemented
}

// GetShow returns a Show
// or errors with ErrNotFound, or ErrInternalServer
func (l *TraktLibrary) GetShow(showID string) (*Show, error) {
	return nil, nil
}

// GetShows returns ErrNotImplemented
func (l *TraktLibrary) GetShows() ([]*Show, error) {
	return nil, ErrNotImplemented
}

// GetSeason returns a Season
// or errors with ErrNotFound, or ErrInternalServer
func (l *TraktLibrary) GetSeason(seasonID string) (*Season, error) {
	return nil, nil
}

// GetSeasonByNumber returns a Season given a Show's ID and a Season number
// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
func (l *TraktLibrary) GetSeasonByNumber(showID string, seasonNumber int) (*Season, error) {
	return nil, nil
}

// GetEpisode returns an Episode
// or errors with ErrNotFound, or ErrInternalServer
func (l *TraktLibrary) GetEpisode(episodeID string) (*Episode, error) {
	return nil, nil
}

// GetEpisodeByNumber returns an Episode  given a Show's ID a Season number
// and Episode's number
// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
func (l *TraktLibrary) GetEpisodeByNumber(showID string, seasonNumber, episodeNumber int) (*Episode, error) {
	return nil, nil
}

// QueryShowsByTitle returns all Shows that match a partial title ordered
// by their probability
// or errors with ErrInternalServer
func (l *TraktLibrary) QueryShowsByTitle(title string) ([]*Show, error) {
	return []*Show{}, nil
}
