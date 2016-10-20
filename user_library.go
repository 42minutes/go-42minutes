package minutes

// UserLibrary holds information on the shows, seasons, and episodes that
// the user has.
type UserLibrary interface {
	// UpsertShow adds or updates a show
	// or error with ErrNotImplemented, or ErrInternalServer
	UpsertShow(show *UserShow) error
	// UpsertSeason adds or updates a season
	// or errors with ErrNotImplemented, or ErrInternalServer, or ErrMissingShow
	UpsertSeason(season *UserSeason) error
	// UpsertEpisode adds or updates a episode
	// or errors with ErrNotImplemented, or ErrInternalServer, ErrMissingShow
	// or ErrMissingSeason
	UpsertEpisode(episode *UserEpisode) error
	// GetShow returns a UserShow
	// or errors with ErrNotFound, or ErrInternalServer
	GetShow(showID string) (*UserShow, error)
	// GetShows returns all Shows
	// or errors with ErrNotImplemented, or ErrInternalServer
	GetShows() ([]*UserShow, error)
	// GetSeasons returns all Seasons for a show
	// or errors with ErrNotFound, or ErrInternalServer
	GetSeasons(showID string) ([]*UserSeason, error)
	// GetSeason returns a UserSeason given a UserShow's ID and a UserSeason number
	// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
	GetSeason(showID string, seasonNumber int) (*UserSeason, error)
	// GetEpisodes returns all Shows for a show and season number
	// or errors with ErrNotFound, or ErrInternalServer
	GetEpisodes(showID string, seasonNumber int) ([]*UserEpisode, error)
	// GetEpisode returns a UserEpisode  given a UserShow's ID a UserSeason number
	// and UserEpisode's number
	// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
	GetEpisode(showID string, seasonNumber, episodeNumber int) (*UserEpisode, error)
	// QueryShowsByTitle returns all Shows that match a partial title ordered
	// by their probability
	// or errors with ErrInternalServer
	QueryShowsByTitle(title string) ([]*UserShow, error)
	// QueryEpisodesForFinder
	QueryEpisodesForFinder() ([]*UserEpisode, error)
	// QueryEpisodesForDownloader
	QueryEpisodesForDownloader() ([]*UserEpisode, error)
}
