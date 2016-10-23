package minutes

import (
	"encoding/json"
	"fmt"
	"time"
)

// UserShow struct
type UserShow struct {
	ID            string `json:"id" gorethink:"id"`
	AiredEpisodes int    `json:"aired_episodes" gorethink:"-" gorm:"-"`
	Airs          struct {
		Day      string `json:"day" gorethink:"-" gorm:"-"`
		Time     string `json:"time" gorethink:"-" gorm:"-"`
		Timezone string `json:"timezone" gorethink:"-" gorm:"-"`
	} `json:"airs" gorethink:"-" gorm:"-"`
	AvailableTranslations []string `json:"-" gorethink:"-" gorm:"-"`
	Certification         string   `json:"-" gorethink:"-" gorm:"-"`
	Country               string   `json:"country" gorethink:"-" gorm:"-"`
	FirstAired            string   `json:"first_aired" gorethink:"-" gorm:"-"`
	Genres                []string `json:"genres" gorethink:"-" gorm:"-"`
	Homepage              string   `json:"-" gorethink:"-" gorm:"-"`
	IDs                   struct {
		Imdb   string `json:"imdb" gorethink:"-" gorm:"-"`
		Slug   string `json:"slug" gorethink:"-" gorm:"-"`
		Tmdb   int    `json:"tmdb" gorethink:"-" gorm:"-"`
		Trakt  int    `json:"trakt" gorethink:"-" gorm:"-"`
		Tvdb   int    `json:"tvdb" gorethink:"-" gorm:"-"`
		Tvrage int    `json:"tvrage" gorethink:"-" gorm:"-"`
	} `json:"ids" gorethink:"-" gorm:"-"`
	Images struct {
		Banner struct {
			Full string `json:"full" gorethink:"-" gorm:"-"`
		} `json:"banner" gorethink:"-" gorm:"-"`
		Clearart struct {
			Full string `json:"full" gorethink:"-" gorm:"-"`
		} `json:"clearart" gorethink:"-" gorm:"-"`
		Fanart struct {
			Full   string `json:"full" gorethink:"-" gorm:"-"`
			Medium string `json:"medium" gorethink:"-" gorm:"-"`
			Thumb  string `json:"thumb" gorethink:"-" gorm:"-"`
		} `json:"fanart" gorethink:"-" gorm:"-"`
		Logo struct {
			Full string `json:"full" gorethink:"-" gorm:"-"`
		} `json:"logo" gorethink:"-" gorm:"-"`
		Poster struct {
			Full   string `json:"full" gorethink:"-" gorm:"-"`
			Medium string `json:"medium" gorethink:"-" gorm:"-"`
			Thumb  string `json:"thumb" gorethink:"-" gorm:"-"`
		} `json:"poster" gorethink:"-" gorm:"-"`
		Thumb struct {
			Full string `json:"full" gorethink:"-" gorm:"-"`
		} `json:"thumb" gorethink:"-" gorm:"-"`
	} `json:"images" gorethink:"-" gorm:"-"`
	Language  string  `json:"language" gorethink:"-" gorm:"-"`
	Network   string  `json:"network" gorethink:"-" gorm:"-"`
	Overview  string  `json:"overview" gorethink:"-" gorm:"-"`
	Rating    float64 `json:"rating" gorethink:"-" gorm:"-"`
	Runtime   float64 `json:"runtime" gorethink:"-" gorm:"-"`
	Status    string  `json:"status" gorethink:"-" gorm:"-"`
	Title     string  `json:"title" gorethink:"title,omitempty"`
	Trailer   string  `json:"-" gorethink:"-" gorm:"-"`
	UpdatedAt string  `json:"-" gorethink:"-" `
	Votes     int     `json:"-" gorethink:"-" gorm:"-"`
	Year      int     `json:"year" gorethink:"-" gorm:"-"`
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
	ShowID       string `json:"show_id" gorethink:"show_id"`
	EpisodeCount int    `json:"episode_count" gorethink:"-" gorm:"-"`
	IDs          struct {
		Tmdb   int `json:"tmdb" gorethink:"-" gorm:"-"`
		Trakt  int `json:"trakt" gorethink:"-" gorm:"-"`
		Tvdb   int `json:"tvdb" gorethink:"-" gorm:"-"`
		Tvrage int `json:"tvrage" gorethink:"-" gorm:"-"`
	} `json:"ids" gorethink:"-" gorm:"-"`
	Images struct {
		Poster struct {
			Full   string `json:"full" gorethink:"-" gorm:"-"`
			Medium string `json:"medium" gorethink:"-" gorm:"-"`
			Thumb  string `json:"thumb" gorethink:"-" gorm:"-"`
		} `json:"poster" gorethink:"-" gorm:"-"`
		Thumb struct {
			Full string `json:"full" gorethink:"-" gorm:"-"`
		} `json:"thumb" gorethink:"-" gorm:"-"`
	} `json:"images" gorethink:"-" gorm:"-"`
	Number   int     `json:"number" gorethink:"number"`
	Overview string  `json:"overview" gorethink:"-" gorm:"-"`
	Rating   float64 `json:"rating" gorethink:"-" gorm:"-"`
	Votes    int     `json:"votes" gorethink:"-" gorm:"-"`

	CID []interface{} `json:"-" gorethink:"id" gorm:"-"`
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
	AvailableTranslations []string   `json:"-" gorethink:"-" gorm:"-"`
	FirstAired            *time.Time `json:"first_aired" gorethink:"-" gorm:"-"`
	IDs                   struct {
		Imdb   string `json:"imdb" gorethink:"-" gorm:"-"`
		Tmdb   int    `json:"tmdb" gorethink:"-" gorm:"-"`
		Trakt  int    `json:"trakt" gorethink:"-" gorm:"-"`
		Tvdb   int    `json:"tvdb" gorethink:"-" gorm:"-"`
		Tvrage int    `json:"tvrage" gorethink:"-" gorm:"-"`
	} `json:"ids" gorethink:"-" gorm:"-"`
	Images struct {
		Screenshot struct {
			Full   string `json:"full" gorethink:"-" gorm:"-"`
			Medium string `json:"medium" gorethink:"-" gorm:"-"`
			Thumb  string `json:"thumb" gorethink:"-" gorm:"-"`
		} `json:"screenshot" gorethink:"-" gorm:"-"`
	} `json:"images" gorethink:"-" gorm:"-"`
	Number    int     `json:"number" gorethink:"number,omitempty"`
	Overview  string  `json:"overview" gorethink:"-" gorm:"-"`
	Rating    float64 `json:"-" gorethink:"-" gorm:"-"`
	Season    int     `json:"season" gorethink:"season,omitempty"`
	Title     string  `json:"title" gorethink:"title,omitempty"`
	UpdatedAt string  `json:"updated_at" gorethink:"-" `
	Votes     int     `json:"-" gorethink:"-" gorm:"-"`

	CID   []interface{} `json:"-" gorethink:"id" gorm:"-"`
	Files []*UserFile   `json:"files" gorethink:"files" gorm:"-"`
}

// UserFile -
type UserFile struct {
	ID           uint   `json:"id" gorm:"primary_key"`
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
	ShowID       string `json:"-" gorethink:"_"`
	Season       int    `json:"-" gorethink:"_"`
	Episode      int    `json:"-" gorethink:"_"`
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
