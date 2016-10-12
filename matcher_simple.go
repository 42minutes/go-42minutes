package minutes

// SimpleMatch uses regular expressions to match against a show, season,
// and episode
type SimpleMatch struct {
	globalLibrary Library
}

// NewSimpleMatch returns a SimpleMatch
func NewSimpleMatch(glib Library) (*SimpleMatch, error) {
	return &SimpleMatch{
		globalLibrary: glib,
	}, nil
}

// Match returns all episodes that match a filename or full path,
// ordered by their probability
// or errors with ErrInternalServer
func (m *SimpleMatch) Match(filename string) ([]*Episode, error) {
	return []*Episode{}, nil
}
