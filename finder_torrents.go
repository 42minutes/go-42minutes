package minutes

// TorrentFinder is an implementation of Finder that can search for torrent
// files and magnet links
type TorrentFinder struct{}

// Find Downloadable Torrents for a given Episode
func (f *TorrentFinder) Find(episode *Episode) ([]Downloadable, error) {
	return []Downloadable{}, nil
}
