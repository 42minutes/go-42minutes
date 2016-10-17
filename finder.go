package minutes

// Finder is responsible for finding Downloadables, caching is also left
// to the implementation and is optional
type Finder interface {
	// Find returns a list of Downloadables for a given Episode
	Find(show *Show, episode *Episode) ([]Downloadable, error)
}
