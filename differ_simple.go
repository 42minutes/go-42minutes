package minutes

// SimpleDiff is used to find missing episodes
type SimpleDiff struct{}

// Diff returns episodes missing from the user's Library
// or returns ErrInternalServer
func (d *SimpleDiff) Diff(user, global *Show) (diff []*Episode, err error) {
	return []*Episode{}, nil
}
