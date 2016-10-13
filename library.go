package minutes

import "errors"

var (
	// ErrNotFound is returned when a single resource was requested
	// but was not found
	ErrNotFound = errors.New("Not found")
	// ErrInternalServer is returned when the request failed with
	// anything that was out of the requester's control
	ErrInternalServer = errors.New("Internal server error")
	// ErrMissingShow is returned when trying to request, add, or modify
	// a resource that is related to a Show, but the Show does not exist
	ErrMissingShow = errors.New("Show does not exist")
	// ErrMissingSeason is returned when trying to request, add, or modify
	// a resource that is related to a Season, but the Season does not exist
	ErrMissingSeason = errors.New("Season does not exist")
	// ErrNotImplemented is returned when trying to add a new resource on a
	// read-only Library
	ErrNotImplemented = errors.New("Not implemented")
)

// Library holds all available information about shows, seasons, and episodes
// By default uses IMDB IDs as unique identifier, a string
// There should be at least two implementation of this interface
// * A global library that should contain all known shows, seasons, and episodes
// * A user-specific library that holds the shows, seasons, and episodes of a
//   single user
type Library interface {
	// UpsertShow adds or updates a show
	// or error with ErrNotImplemented, or ErrInternalServer
	UpsertShow(show *Show) error
	// UpsertSeason adds or updates a season
	// or errors with ErrNotImplemented, or ErrInternalServer, or ErrMissingShow
	UpsertSeason(season *Season) error
	// UpsertEpisode adds or updates a episode
	// or errors with ErrNotImplemented, or ErrInternalServer, ErrMissingShow
	// or ErrMissingSeason
	UpsertEpisode(episode *Episode) error
	// GetShow returns a Show
	// or errors with ErrNotFound, or ErrInternalServer
	GetShow(showID string) (*Show, error)
	// GetShows returns all Shows
	// or errors with ErrNotImplemented, or ErrInternalServer
	GetShows() ([]*Show, error)
	// GetSeason returns a Season
	// or errors with ErrNotFound, or ErrInternalServer
	GetSeason(seasonID string) (*Season, error)
	// GetSeasonsByShow returns all Seasons for a show
	// or errors with ErrNotFound, or ErrInternalServer
	GetSeasonsByShow(showID string) ([]*Season, error)
	// GetSeasonByNumber returns a Season given a Show's ID and a Season number
	// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
	GetSeasonByNumber(showID string, seasonNumber int) (*Season, error)
	// GetEpisode returns an Episode
	// or errors with ErrNotFound, or ErrInternalServer
	GetEpisode(episodeID string) (*Episode, error)
	// GetEpisodesBySeason returns all Episodes for a Season
	// or errors with ErrNotFound, or ErrInternalServer
	GetEpisodesBySeason(seasonID string) ([]*Episode, error)
	// GetEpisodesBySeasonNumber returns all Shows for a show and season number
	// or errors with ErrNotFound, or ErrInternalServer
	GetEpisodesBySeasonNumber(showID string, seasonNumber int) ([]*Episode, error)
	// GetEpisodeByNumber returns an Episode  given a Show's ID a Season number
	// and Episode's number
	// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
	GetEpisodeByNumber(showID string, seasonNumber, episodeNumber int) (*Episode, error)
	// QueryShowsByTitle returns all Shows that match a partial title ordered
	// by their probability
	// or errors with ErrInternalServer
	QueryShowsByTitle(title string) ([]*Show, error)
}
