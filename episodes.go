package minutes

import "time"

// Show struct according to the Trakt v2 API
type Show struct {
	ID            string `json:"id"`
	AiredEpisodes int    `json:"aired_episodes"`
	Airs          struct {
		Day      string `json:"day"`
		Time     string `json:"time"`
		Timezone string `json:"timezone"`
	} `json:"airs"`
	AvailableTranslations []string `json:"available_translations"`
	Certification         string   `json:"certification"`
	Country               string   `json:"country"`
	FirstAired            string   `json:"first_aired"`
	Genres                []string `json:"genres"`
	Homepage              string   `json:"homepage"`
	IDs                   struct {
		Imdb   string `json:"imdb"`
		Slug   string `json:"slug"`
		Tmdb   int    `json:"tmdb"`
		Trakt  int    `json:"trakt"`
		Tvdb   int    `json:"tvdb"`
		Tvrage int    `json:"tvrage"`
	} `json:"ids"`
	Images struct {
		Banner struct {
			Full string `json:"full"`
		} `json:"banner"`
		Clearart struct {
			Full string `json:"full"`
		} `json:"clearart"`
		Fanart struct {
			Full   string `json:"full"`
			Medium string `json:"medium"`
			Thumb  string `json:"thumb"`
		} `json:"fanart"`
		Logo struct {
			Full string `json:"full"`
		} `json:"logo"`
		Poster struct {
			Full   string `json:"full"`
			Medium string `json:"medium"`
			Thumb  string `json:"thumb"`
		} `json:"poster"`
		Thumb struct {
			Full string `json:"full"`
		} `json:"thumb"`
	} `json:"images"`
	Language  string  `json:"language"`
	Network   string  `json:"network"`
	Overview  string  `json:"overview"`
	Rating    float64 `json:"rating"`
	Runtime   float64 `json:"runtime"`
	Status    string  `json:"status"`
	Title     string  `json:"title"`
	Trailer   string  `json:"trailer"`
	UpdatedAt string  `json:"updated_at"`
	Votes     int     `json:"votes"`
	Year      int     `json:"year"`
}

// Season struct according to the Trakt v2 API
type Season struct {
	ShowID       string `json:"show_id"`
	EpisodeCount int    `json:"episode_count"`
	IDs          struct {
		Tmdb   int `json:"tmdb"`
		Trakt  int `json:"trakt"`
		Tvdb   int `json:"tvdb"`
		Tvrage int `json:"tvrage"`
	} `json:"ids"`
	Images struct {
		Poster struct {
			Full   string `json:"full"`
			Medium string `json:"medium"`
			Thumb  string `json:"thumb"`
		} `json:"poster"`
		Thumb struct {
			Full string `json:"full"`
		} `json:"thumb"`
	} `json:"images"`
	Number   int     `json:"number"`
	Overview string  `json:"overview"`
	Rating   float64 `json:"rating"`
	Votes    int     `json:"votes"`
}

// Episode struct according to the Trakt v2 API
type Episode struct {
	ShowID                string     `json:"show_id"`
	SeasonID              string     `json:"season_id"`
	AvailableTranslations []string   `json:"available_translations"`
	FirstAired            *time.Time `json:"first_aired"`
	IDs                   struct {
		Imdb   string `json:"imdb"`
		Tmdb   int    `json:"tmdb"`
		Trakt  int    `json:"trakt"`
		Tvdb   int    `json:"tvdb"`
		Tvrage int    `json:"tvrage"`
	} `json:"ids"`
	Images struct {
		Screenshot struct {
			Full   string `json:"full"`
			Medium string `json:"medium"`
			Thumb  string `json:"thumb"`
		} `json:"screenshot"`
	} `json:"images"`
	Number    int     `json:"number"`
	Overview  string  `json:"overview"`
	Rating    float64 `json:"rating"`
	Season    int     `json:"season"`
	Title     string  `json:"title"`
	UpdatedAt string  `json:"updated_at"`
	Votes     int     `json:"votes"`
}
