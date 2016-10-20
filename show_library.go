package minutes

// ShowLibrary holds all available information about shows, seasons, and episodes
// By default uses IMDB IDs as unique identifier, a string
type ShowLibrary interface {
	// GetShow returns a Show
	// or errors with ErrNotFound, or ErrInternalServer
	GetShow(showID string) (*Show, error)
	// GetSeasons returns all Seasons for a show
	// or errors with ErrNotFound, or ErrInternalServer
	GetSeasons(showID string) ([]*Season, error)
	// GetSeason returns a Season given a Show's ID and a Season number
	// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
	GetSeason(showID string, seasonNumber int) (*Season, error)
	// GetEpisodes returns all Shows for a show and season number
	// or errors with ErrNotFound, or ErrInternalServer
	GetEpisodes(showID string, seasonNumber int) ([]*Episode, error)
	// GetEpisode returns an Episode  given a Show's ID a Season number
	// and Episode's number
	// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
	GetEpisode(showID string, seasonNumber, episodeNumber int) (*Episode, error)
	// QueryShowsByTitle returns all Shows that match a partial title ordered
	// by their probability
	// or errors with ErrInternalServer
	QueryShowsByTitle(title string) ([]*Show, error)
}
