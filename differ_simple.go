package minutes

// SimpleDiff is used to find missing episodes
type SimpleDiff struct{}

// Diff returns episodes missing from the user's Library
// or returns ErrInternalServer
func (d *SimpleDiff) Diff(user, global Library) (diff []*Episode, err error) {
	result := []*Episode{}

	uShows, err := user.GetShows()
	if err != nil {
		return nil, ErrInternalServer
	}

	// Parse all user shows
	for _, ush := range uShows {
		// Find show in global library
		globalSeasons, err := global.GetSeasonsByShow(ush.ID)
		if err != nil {
			return nil, ErrInternalServer
		}

		for _, gseas := range globalSeasons {
			// Get all episodes from global
			gEpisodes, err := global.GetEpisodesBySeasonNumber(ush.ID, gseas.Number)
			if err != nil {
				return nil, ErrInternalServer
			}

			// For each episode try to match in local
			for _, gepp := range gEpisodes {
				// Try to find episode in user lib
				_, err := user.GetEpisodeByNumber(ush.ID, gseas.Number, gepp.Number)
				switch err {
				case ErrNotFound:
					result = append(result, gepp)
				case ErrInternalServer:
					return nil, ErrInternalServer
				}
			}
		}

	}

	return result, nil
}
