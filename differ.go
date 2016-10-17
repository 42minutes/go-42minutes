package minutes

// Differ is used to find missing episodes
type Differ interface {
	// Diff returns episodes missing from the user's Library
	// or returns ErrInternalServer
	Diff(user, global Library) (diff []*Episode, err error)
}
