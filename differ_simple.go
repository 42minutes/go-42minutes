package minutes

// SimpleDiff is used to find missing episodes
type SimpleDiff struct {
	ulibrary Library
	glibrary Library
}

func NewSimpleDiff(ulib, glib Library) *SimpleDiff {
	return &SimpleDiff{
		ulibrary: ulib,
		glibrary: glib,
	}
}

// Diff returns episodes missing from the user's Library
// or returns ErrInternalServer
func (d *SimpleDiff) Diff(ush, gsh *Show) (diff []*Episode, err error) {
	result := []*Episode{}

	// Find show in global library
	globalSeasons, err := d.glibrary.GetSeasonsByShow(ush.ID)
	if err != nil {
		return nil, ErrInternalServer
	}

	for _, gseas := range globalSeasons {
		// Get all episodes from global
		gEpisodes, err := d.glibrary.GetEpisodesBySeasonNumber(ush.ID, gseas.Number)
		if err != nil {
			return nil, ErrInternalServer
		}

		// For each episode try to match in local
		for _, gepp := range gEpisodes {
			// Try to find episode in user lib
			_, err := d.ulibrary.GetEpisodeByNumber(ush.ID, gseas.Number, gepp.Number)
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
