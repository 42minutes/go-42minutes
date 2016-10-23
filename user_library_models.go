package minutes

import (
	"encoding/json"
	"fmt"
	"time"
)

// UserShow struct
type UserShow struct {
	ID            string `json:"id" gorethink:"id"`
	AiredEpisodes int    `json:"aired_episodes" gorethink:"-"`
	Airs          struct {
		Day      string `json:"day" gorethink:"-"`
		Time     string `json:"time" gorethink:"-"`
		Timezone string `json:"timezone" gorethink:"-"`
	} `json:"airs" gorethink:"-"`
	AvailableTranslations []string `json:"-" gorethink:"-"`
	Certification         string   `json:"-" gorethink:"-"`
	Country               string   `json:"country" gorethink:"-"`
	FirstAired            string   `json:"first_aired" gorethink:"-"`
	Genres                []string `json:"genres" gorethink:"-"`
	Homepage              string   `json:"-" gorethink:"-"`
	IDs                   struct {
		Imdb   string `json:"imdb" gorethink:"-"`
		Slug   string `json:"slug" gorethink:"-"`
		Tmdb   int    `json:"tmdb" gorethink:"-"`
		Trakt  int    `json:"trakt" gorethink:"-"`
		Tvdb   int    `json:"tvdb" gorethink:"-"`
		Tvrage int    `json:"tvrage" gorethink:"-"`
	} `json:"ids" gorethink:"-"`
	Images struct {
		Banner struct {
			Full string `json:"full" gorethink:"-"`
		} `json:"banner" gorethink:"-"`
		Clearart struct {
			Full string `json:"full" gorethink:"-"`
		} `json:"clearart" gorethink:"-"`
		Fanart struct {
			Full   string `json:"full" gorethink:"-"`
			Medium string `json:"medium" gorethink:"-"`
			Thumb  string `json:"thumb" gorethink:"-"`
		} `json:"fanart" gorethink:"-"`
		Logo struct {
			Full string `json:"full" gorethink:"-"`
		} `json:"logo" gorethink:"-"`
		Poster struct {
			Full   string `json:"full" gorethink:"-"`
			Medium string `json:"medium" gorethink:"-"`
			Thumb  string `json:"thumb" gorethink:"-"`
		} `json:"poster" gorethink:"-"`
		Thumb struct {
			Full string `json:"full" gorethink:"-"`
		} `json:"thumb" gorethink:"-"`
	} `json:"images" gorethink:"-"`
	Language  string  `json:"language" gorethink:"-"`
	Network   string  `json:"network" gorethink:"-"`
	Overview  string  `json:"overview" gorethink:"-"`
	Rating    float64 `json:"rating" gorethink:"-"`
	Runtime   float64 `json:"runtime" gorethink:"-"`
	Status    string  `json:"status" gorethink:"-"`
	Title     string  `json:"title" gorethink:"title,omitempty"`
	Trailer   string  `json:"-" gorethink:"-"`
	UpdatedAt string  `json:"-" gorethink:"-"`
	Votes     int     `json:"-" gorethink:"-"`
	Year      int     `json:"year" gorethink:"-"`
}

// MergeInPlace -
func (sh *UserShow) MergeInPlace(gsh *Show) {
	data, _ := json.Marshal(gsh)
	json.Unmarshal(data, &sh)
}

// GetCID -
func (sh *UserShow) GetCID() interface{} {
	return sh.ID
}

// UserSeason struct
type UserSeason struct {
	ShowID       string `json:"show_id" gorethink:"show_id" gorm:"-"`
	EpisodeCount int    `json:"episode_count" gorethink:"-"`
	IDs          struct {
		Tmdb   int `json:"tmdb" gorethink:"-"`
		Trakt  int `json:"trakt" gorethink:"-"`
		Tvdb   int `json:"tvdb" gorethink:"-"`
		Tvrage int `json:"tvrage" gorethink:"-"`
	} `json:"ids" gorethink:"-"`
	Images struct {
		Poster struct {
			Full   string `json:"full" gorethink:"-"`
			Medium string `json:"medium" gorethink:"-"`
			Thumb  string `json:"thumb" gorethink:"-"`
		} `json:"poster" gorethink:"-"`
		Thumb struct {
			Full string `json:"full" gorethink:"-"`
		} `json:"thumb" gorethink:"-"`
	} `json:"images" gorethink:"-"`
	Number   int     `json:"number" gorethink:"number"`
	Overview string  `json:"overview" gorethink:"-"`
	Rating   float64 `json:"rating" gorethink:"-"`
	Votes    int     `json:"votes" gorethink:"-"`

	CID []interface{} `json:"-" gorethink:"id"`
}

// MergeInPlace -
func (se *UserSeason) MergeInPlace(gse *Season) {
	data, _ := json.Marshal(gse)
	json.Unmarshal(data, &se)
}

// GetCID -
func (se *UserSeason) GetCID() []interface{} {
	return []interface{}{
		se.ShowID,
		se.Number,
	}
}

// UserEpisode struct
type UserEpisode struct {
	ShowID                string     `json:"show_id" gorethink:"show_id"`
	AvailableTranslations []string   `json:"-" gorethink:"-"`
	FirstAired            *time.Time `json:"first_aired" gorethink:"-"`
	IDs                   struct {
		Imdb   string `json:"imdb" gorethink:"-"`
		Tmdb   int    `json:"tmdb" gorethink:"-"`
		Trakt  int    `json:"trakt" gorethink:"-"`
		Tvdb   int    `json:"tvdb" gorethink:"-"`
		Tvrage int    `json:"tvrage" gorethink:"-"`
	} `json:"ids" gorethink:"-"`
	Images struct {
		Screenshot struct {
			Full   string `json:"full" gorethink:"-"`
			Medium string `json:"medium" gorethink:"-"`
			Thumb  string `json:"thumb" gorethink:"-"`
		} `json:"screenshot" gorethink:"-"`
	} `json:"images" gorethink:"-"`
	Number    int     `json:"number" gorethink:"number,omitempty"`
	Overview  string  `json:"overview" gorethink:"-"`
	Rating    float64 `json:"-" gorethink:"-"`
	Season    int     `json:"season" gorethink:"season,omitempty"`
	Title     string  `json:"title" gorethink:"title,omitempty"`
	UpdatedAt string  `json:"updated_at" gorethink:"-"`
	Votes     int     `json:"-" gorethink:"-"`

	CID   []interface{} `json:"-" gorethink:"id" gorm:"-"`
	Files []*UserFile   `json:"files" gorethink:"files"`
}

// UserFile -
type UserFile struct {
	Name         string `json:"name" gorethink:"name"`
	Path         string `json:"path" gorethink:"path"`
	Infohash     string `json:"-" gorethink:"infohash"`
	Resolution   string `json:"resolution" gorethink:"resolution"`
	Source       string `json:"source" gorethink:"source"`
	VideoCodec   string `json:"video_codec" gorethink:"video_codec"`
	AudioCodec   string `json:"audio_codec" gorethink:"audio_codec"`
	ReleaseGroup string `json:"release_group" gorethink:"release_group"`
	CRC32        string `json:"crc32" gorethink:"crc32"`
	Status       string `json:"status" gorethink:"status"`
	RetryTime    int64  `json:"retry_time" gorethink:"retry_time"`
}

// MergeInPlace -
func (ep *UserEpisode) MergeInPlace(gep *Episode) {
	data, _ := json.Marshal(gep)
	json.Unmarshal(data, &ep)
}

// String -
func (ep *UserEpisode) String() string {
	return fmt.Sprintf("ID:%s S%02dE%02d", ep.ShowID, ep.Season, ep.Number)
}

// GetCID -
func (ep *UserEpisode) GetCID() []interface{} {
	return []interface{}{
		ep.ShowID,
		ep.Season,
		ep.Number,
	}
}
