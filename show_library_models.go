package minutes

import "time"

// Show struct according to the Trakt v2 API
type Show struct {
	ID            string `json:"id" gorethink:"id"`
	AiredEpisodes int    `json:"aired_episodes" gorethink:"aired_episodes"`
	Airs          struct {
		Day      string `json:"day" gorethink:"day,omitempty"`
		Time     string `json:"time" gorethink:"time,omitempty"`
		Timezone string `json:"timezone" gorethink:"timezone,omitempty"`
	} `json:"airs" gorethink:"airs,omitempty"`
	AvailableTranslations []string `json:"available_translations" gorethink:"available_translations,omitempty"`
	Certification         string   `json:"certification" gorethink:"certification,omitempty"`
	Country               string   `json:"country" gorethink:"country,omitempty"`
	FirstAired            string   `json:"first_aired" gorethink:"first_aired,omitempty"`
	Genres                []string `json:"genres" gorethink:"genres,omitempty"`
	Homepage              string   `json:"homepage" gorethink:"homepage,omitempty"`
	IDs                   struct {
		Imdb   string `json:"imdb" gorethink:"imdb,omitempty"`
		Slug   string `json:"slug" gorethink:"slug,omitempty"`
		Tmdb   int    `json:"tmdb" gorethink:"tmdb,omitempty"`
		Trakt  int    `json:"trakt" gorethink:"trakt,omitempty"`
		Tvdb   int    `json:"tvdb" gorethink:"tvdb,omitempty"`
		Tvrage int    `json:"tvrage" gorethink:"tvrage,omitempty"`
	} `json:"ids" gorethink:"ids"`
	Images struct {
		Banner struct {
			Full string `json:"full" gorethink:"full,omitempty"`
		} `json:"banner" gorethink:"banner,omitempty"`
		Clearart struct {
			Full string `json:"full" gorethink:"full,omitempty"`
		} `json:"clearart" gorethink:"full,omitempty"`
		Fanart struct {
			Full   string `json:"full" gorethink:"full,omitempty"`
			Medium string `json:"medium" gorethink:"medium,omitempty"`
			Thumb  string `json:"thumb" gorethink:"thumb,omitempty"`
		} `json:"fanart" gorethink:"fanart,omitempty"`
		Logo struct {
			Full string `json:"full" gorethink:"full,omitempty"`
		} `json:"logo" gorethink:"logo,omitempty"`
		Poster struct {
			Full   string `json:"full" gorethink:"full,omitempty"`
			Medium string `json:"medium" gorethink:"medium,omitempty"`
			Thumb  string `json:"thumb" gorethink:"thumb,omitempty"`
		} `json:"poster" gorethink:"poster,omitempty"`
		Thumb struct {
			Full string `json:"full" gorethink:"full,omitempty"`
		} `json:"thumb" gorethink:"thumb,omitempty"`
	} `json:"images" gorethink:"images,omitempty"`
	Language  string  `json:"language" gorethink:"language,omitempty"`
	Network   string  `json:"network" gorethink:"network,omitempty"`
	Overview  string  `json:"overview" gorethink:"overview,omitempty"`
	Rating    float64 `json:"rating" gorethink:"rating,omitempty"`
	Runtime   float64 `json:"runtime" gorethink:"runtime,omitempty"`
	Status    string  `json:"status" gorethink:"status,omitempty"`
	Title     string  `json:"title" gorethink:"title,omitempty"`
	Trailer   string  `json:"trailer" gorethink:"trailer,omitempty"`
	UpdatedAt string  `json:"updated_at" gorethink:"updated_at,omitempty"`
	Votes     int     `json:"votes" gorethink:"votes,omitempty"`
	Year      int     `json:"year" gorethink:"year,omitempty"`
}

// Season struct according to the Trakt v2 API
type Season struct {
	ID           string `json:"id" gorethink:"id"`
	ShowID       string `json:"show_id" gorethink:"show_id"`
	EpisodeCount int    `json:"episode_count" gorethink:"episode_count"`
	IDs          struct {
		Tmdb   int `json:"tmdb" gorethink:"tmdb,omitempty"`
		Trakt  int `json:"trakt" gorethink:"trakt,omitempty"`
		Tvdb   int `json:"tvdb" gorethink:"tvdb,omitempty"`
		Tvrage int `json:"tvrage" gorethink:"tvrage,omitempty"`
	} `json:"ids" gorethink:"ids"`
	Images struct {
		Poster struct {
			Full   string `json:"full" gorethink:"full,omitempty"`
			Medium string `json:"medium" gorethink:"medium,omitempty"`
			Thumb  string `json:"thumb" gorethink:"thumb,omitempty"`
		} `json:"poster" gorethink:"poster,omitempty"`
		Thumb struct {
			Full string `json:"full" gorethink:"full,omitempty"`
		} `json:"thumb" gorethink:"thumb,omitempty"`
	} `json:"images" gorethink:"images,omitempty"`
	Number   int     `json:"number" gorethink:"number,omitempty"`
	Overview string  `json:"overview" gorethink:"overview,omitempty"`
	Rating   float64 `json:"rating" gorethink:"rating,omitempty"`
	Votes    int     `json:"votes" gorethink:"votes,omitempty"`
}

// Episode struct according to the Trakt v2 API
type Episode struct {
	ID                    string     `json:"id" gorethink:"id"`
	ShowID                string     `json:"show_id" gorethink:"show_id"`
	SeasonID              string     `json:"season_id" gorethink:"season_id,omitempty"`
	AvailableTranslations []string   `json:"available_translations" gorethink:"available_translations,omitempty"`
	FirstAired            *time.Time `json:"first_aired" gorethink:"first_aired,omitempty"`
	IDs                   struct {
		Imdb   string `json:"imdb" gorethink:"imdb,omitempty"`
		Tmdb   int    `json:"tmdb" gorethink:"tmdb,omitempty"`
		Trakt  int    `json:"trakt" gorethink:"trakt,omitempty"`
		Tvdb   int    `json:"tvdb" gorethink:"tvdb,omitempty"`
		Tvrage int    `json:"tvrage" gorethink:"tvrage,omitempty"`
	} `json:"ids" gorethink:"ids"`
	Images struct {
		Screenshot struct {
			Full   string `json:"full" gorethink:"full,omitempty"`
			Medium string `json:"medium" gorethink:"medium,omitempty"`
			Thumb  string `json:"thumb" gorethink:"thumb,omitempty"`
		} `json:"screenshot" gorethink:"screenshot,omitempty"`
	} `json:"images" gorethink:"images,omitempty"`
	Number    int     `json:"number" gorethink:"number,omitempty"`
	Overview  string  `json:"overview" gorethink:"overview,omitempty"`
	Rating    float64 `json:"rating" gorethink:"rating,omitempty"`
	Season    int     `json:"season" gorethink:"season,omitempty"`
	Title     string  `json:"title" gorethink:"title,omitempty"`
	UpdatedAt string  `json:"updated_at" gorethink:"updated_at,omitempty"`
	Votes     int     `json:"votes" gorethink:"votes,omitempty"`
}
