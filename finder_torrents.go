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
	_, ih, _ := torrentlookup.ProviderTPB.Search(fmt.Sprintf("%s s%0de%0d", sh.Title, ep.Season, ep.Number))
	// TODO(geoah) Check nm agains matcher maybe
	dl := &DownloadableMagnet{
		Infohash: ih,
		Magnet:   torrentlookup.CreateFakeMagnet(ih),
	}
	return []Downloadable{dl}, nil
}
