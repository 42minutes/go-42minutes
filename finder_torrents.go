package minutes

import (
	"fmt"

	torrentlookup "github.com/42minutes/go-torrentlookup"
)

// TorrentFinder is an implementation of Finder that can search for torrent
// files and magnet links
type TorrentFinder struct{}

// Find Downloadable Torrents for a given Episode
func (f *TorrentFinder) Find(sh *Show, ep *Episode) ([]Downloadable, error) {
	// TODO(geoah) Handle error
	qr := fmt.Sprintf("%s s%02de%02d", sh.Title, ep.Season, ep.Number)
	_, ih, _ := torrentlookup.ProviderTPB.Search(qr)
	// TODO(geoah) Check nm agains matcher maybe
	if ih != "" {
		dl := &MagnetDownloadable{
			Infohash: ih,
			Magnet:   torrentlookup.CreateFakeMagnet(ih),
		}
		return []Downloadable{dl}, nil
	}
	return []Downloadable{}, ErrNotFound
}
