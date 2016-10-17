package minutes

// DownloadableTorrent is an implementation of the Downloadable for torrent files
type DownloadableTorrent struct{}

// GetID returns the torrent's hash
func (d *DownloadableTorrent) GetID() string {
	return ""
}

// DownloadableMagnet is an implementation of the Downloadable for magnet links
type DownloadableMagnet struct {
	Infohash string
	Magnet   string
}

// GetID returns the torrent's hash
func (d *DownloadableMagnet) GetID() string {
	return d.Infohash
}

// DownloaderTorrent is an implementation of the Downloader specifically
// for torrent files and magnet links
type DownloaderTorrent struct{}

// Download adds a Downloadable to the list of things to download
// or errors with ErrDownloadableNotSupported, ErrDownloadableNotComplete,
// or ErrInternalServer
func (d *DownloaderTorrent) Download(Downloadable) error {
	return nil
}

// List returns all Downloadables
// or errors with ErrInternalServer
func (d *DownloaderTorrent) List() ([]Downloadable, error) {
	return []Downloadable{}, nil
}

// Start starts a download
// or errors with ErrNotFound, or ErrInternalServer
func (d *DownloaderTorrent) Start(dID string) error {
	return nil
}

// Stop stops a download
// or errors with ErrNotFound, or ErrInternalServer
func (d *DownloaderTorrent) Stop(dID string) error {
	return nil
}

// Progress returns the Downloadable's progress (%)
// or errors with ErrNotFound, or ErrInternalServer
func (d *DownloaderTorrent) Progress(dID string) error {
	return nil
}
