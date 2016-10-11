package minutes

import "errors"

var (
	// ErrDownloadableNotSupported is returned when the Downloader does not
	// support a type of Downloadable that was added to it
	ErrDownloadableNotSupported = errors.New("Downloadable not supported")
	// ErrDownloadableNotComplete is returned when the Downloader cannot
	// actually download a Downloadable due to missing information
	ErrDownloadableNotComplete = errors.New("Downloadable not complete")
)

// Downloader handles downloading files, either on its own or through the use
// of external apps and services
type Downloader interface {
	// Download adds a Downloadable to the list of things to download
	// or errors with ErrDownloadableNotSupported, ErrDownloadableNotComplete,
	// or ErrInternalServer
	Download(Downloadable) error
	// List returns all Downloadables
	// or errors with ErrInternalServer
	List() ([]Downloadable, error)
	// Start starts a download
	// or errors with ErrNotFound, or ErrInternalServer
	Start(dID string) error
	// Stop stops a download
	// or errors with ErrNotFound, or ErrInternalServer
	Stop(dID string) error
	// Progress returns the Downloadable's progress (%)
	// or errors with ErrNotFound, or ErrInternalServer
	Progress(dID string) error
}

// Downloadable describes anything that a Downloader could download
// Could be a Torrent, NBZ, HTTP, anything
// The Downloader will have to cast the added Downloadables to get any data
// it requires to actually download it
type Downloadable interface {
	// GetID returns a unique identifier for this Downloadable
	GetID() string
}
