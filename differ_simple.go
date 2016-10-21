package minutes

// SimpleDiff is used to find missing episodes
type SimpleDiff struct {
	ulibrary UserLibrary
	glibrary ShowLibrary
}

// NewSimpleDiff -
func NewSimpleDiff(ulib UserLibrary, glib ShowLibrary) *SimpleDiff {
	return &SimpleDiff{
		ulibrary: ulib,
		glibrary: glib,
	}
}

// Diff returns episodes missing from the user's Library
// or returns ErrInternalServer
func (d *SimpleDiff) Diff(ush *UserShow, gsh *Show) (diff []*Episode, err error) {
	result := []*Episode{}

	// Find show in global library
	globalSeasons, err := d.glibrary.GetSeasons(ush.ID)
	if err != nil {
		return nil, ErrInternalServer
	}

	for _, gseas := range globalSeasons {
		// TODO(geoah) At some point we might care about specials
		if gseas.Number == 0 {
			continue
		}
		// Get all episodes from global
		gEpisodes, err := d.glibrary.GetEpisodes(ush.ID, gseas.Number)
		if err != nil {
			return nil, ErrInternalServer
		}

		// For each episode try to match in local
		for _, gepp := range gEpisodes {
			// Try to find episode in user lib
			_, err := d.ulibrary.GetEpisode(ush.ID, gseas.Number, gepp.Number)
			switch err {
			case ErrNotFound:
				result = append(result, gepp)
			case ErrInternalServer:
				return nil, ErrInternalServer
			}
		}
	}

	return result, nil
}
