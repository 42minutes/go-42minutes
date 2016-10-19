package minutes

import (
	"fmt"
	"io/ioutil"
)

// DownloadableTorrent is an implementation of the Downloadable for torrent files
type DownloadableTorrent struct{}

// GetID returns the torrent's hash
func (d *DownloadableTorrent) GetID() string {
	return ""
}

// MagnetDownloadable is an implementation of the Downloadable for magnet links
type MagnetDownloadable struct {
	Infohash string
	Magnet   string
}

// GetID returns the torrent's hash
func (d *MagnetDownloadable) GetID() string {
	return d.Infohash
}

// TorrentDownloader is an implementation of the Downloader specifically
// for torrent files and magnet links.
// The first version of this downloader only gets get .torrent and .magent
// metadata-only files for other applications to download.
type TorrentDownloader struct {
	destination string
}

// NewTorrentDownloader returns a TorrentDownloader given a destination path
func NewTorrentDownloader(dst string) *TorrentDownloader {
	return &TorrentDownloader{
		destination: dst,
	}
}

// Download adds a Downloadable to the list of things to download
// or errors with ErrDownloadableNotSupported, ErrDownloadableNotComplete,
// or ErrInternalServer
func (d *TorrentDownloader) Download(dnl Downloadable) error {
	switch v := dnl.(type) {
	case *MagnetDownloadable:
		fn := fmt.Sprintf("%s/%s.magnet", d.destination, dnl.GetID())
		if err := ioutil.WriteFile(fn, []byte(v.Magnet), 0644); err != nil {
			return err
		}
		return nil
	default:
		return ErrNotImplemented
	}
}
