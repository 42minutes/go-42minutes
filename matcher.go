package minutes

// Matcher tries to match a filename with an episode
type Matcher interface {
	// Match returns all episodes that match a filename or full path,
	// ordered by their probability
	// or errors with ErrInternalServer
	Match(filename string) ([]*Episode, error)
}
