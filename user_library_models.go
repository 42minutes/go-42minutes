package minutes

import (
	"encoding/json"
	"fmt"
	"time"
)

// UserShow struct
type UserShow struct {
	ID            string `json:"id" gorethink:"id" gorm:"primary_key"`
	AiredEpisodes int    `json:"aired_episodes" gorethink:"-" sql:"-"`
	Airs          struct {
		Day      string `json:"day" gorethink:"-" sql:"-"`
		Time     string `json:"time" gorethink:"-" sql:"-"`
		Timezone string `json:"timezone" gorethink:"-" sql:"-"`
	} `json:"airs" gorethink:"-" sql:"-"`
	AvailableTranslations []string `json:"-" gorethink:"-" sql:"-"`
	Certification         string   `json:"-" gorethink:"-" sql:"-"`
	Country               string   `json:"country" gorethink:"-" sql:"-"`
	FirstAired            string   `json:"first_aired" gorethink:"-" sql:"-"`
	Genres                []string `json:"genres" gorethink:"-" sql:"-"`
	Homepage              string   `json:"-" gorethink:"-" sql:"-"`
	IDs                   struct {
		Imdb   string `json:"imdb" gorethink:"-" sql:"-"`
		Slug   string `json:"slug" gorethink:"-" sql:"-"`
		Tmdb   int    `json:"tmdb" gorethink:"-" sql:"-"`
		Trakt  int    `json:"trakt" gorethink:"-" sql:"-"`
		Tvdb   int    `json:"tvdb" gorethink:"-" sql:"-"`
		Tvrage int    `json:"tvrage" gorethink:"-" sql:"-"`
	} `json:"ids" gorethink:"-" sql:"-"`
	Images struct {
		Banner struct {
			Full string `json:"full" gorethink:"-" sql:"-"`
		} `json:"banner" gorethink:"-" sql:"-"`
		Clearart struct {
			Full string `json:"full" gorethink:"-" sql:"-"`
		} `json:"clearart" gorethink:"-" sql:"-"`
		Fanart struct {
			Full   string `json:"full" gorethink:"-" sql:"-"`
			Medium string `json:"medium" gorethink:"-" sql:"-"`
			Thumb  string `json:"thumb" gorethink:"-" sql:"-"`
		} `json:"fanart" gorethink:"-" sql:"-"`
		Logo struct {
			Full string `json:"full" gorethink:"-" sql:"-"`
		} `json:"logo" gorethink:"-" sql:"-"`
		Poster struct {
			Full   string `json:"full" gorethink:"-" sql:"-"`
			Medium string `json:"medium" gorethink:"-" sql:"-"`
			Thumb  string `json:"thumb" gorethink:"-" sql:"-"`
		} `json:"poster" gorethink:"-" sql:"-"`
		Thumb struct {
			Full string `json:"full" gorethink:"-" sql:"-"`
		} `json:"thumb" gorethink:"-" sql:"-"`
	} `json:"images" gorethink:"-" sql:"-"`
	Language  string  `json:"language" gorethink:"-" sql:"-"`
	Network   string  `json:"network" gorethink:"-" sql:"-"`
	Overview  string  `json:"overview" gorethink:"-" sql:"-"`
	Rating    float64 `json:"rating" gorethink:"-" sql:"-"`
	Runtime   float64 `json:"runtime" gorethink:"-" sql:"-"`
	Status    string  `json:"status" gorethink:"-" sql:"-"`
	Title     string  `json:"title" gorethink:"title,omitempty"`
	Trailer   string  `json:"-" gorethink:"-" sql:"-"`
	UpdatedAt string  `json:"-" gorethink:"-" `
	Votes     int     `json:"-" gorethink:"-" sql:"-"`
	Year      int     `json:"year" gorethink:"-" sql:"-"`
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
	ID           uint   `json:"id" gorm:"primary_key"`
	ShowID       string `json:"show_id" gorethink:"show_id"`
	EpisodeCount int    `json:"episode_count" gorethink:"-" sql:"-"`
	IDs          struct {
		Tmdb   int `json:"tmdb" gorethink:"-" sql:"-"`
		Trakt  int `json:"trakt" gorethink:"-" sql:"-"`
		Tvdb   int `json:"tvdb" gorethink:"-" sql:"-"`
		Tvrage int `json:"tvrage" gorethink:"-" sql:"-"`
	} `json:"ids" gorethink:"-" sql:"-"`
	Images struct {
		Poster struct {
			Full   string `json:"full" gorethink:"-" sql:"-"`
			Medium string `json:"medium" gorethink:"-" sql:"-"`
			Thumb  string `json:"thumb" gorethink:"-" sql:"-"`
		} `json:"poster" gorethink:"-" sql:"-"`
		Thumb struct {
			Full string `json:"full" gorethink:"-" sql:"-"`
		} `json:"thumb" gorethink:"-" sql:"-"`
	} `json:"images" gorethink:"-" sql:"-"`
	Number   int     `json:"number" gorethink:"number"`
	Overview string  `json:"overview" gorethink:"-" sql:"-"`
	Rating   float64 `json:"rating" gorethink:"-" sql:"-"`
	Votes    int     `json:"votes" gorethink:"-" sql:"-"`

	CID []interface{} `json:"-" gorethink:"id" sql:"-"`
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
	ID                    uint       `json:"id" gorm:"primary_key"`
	ShowID                string     `json:"show_id" gorethink:"show_id"`
	AvailableTranslations []string   `json:"-" gorethink:"-" sql:"-"`
	FirstAired            *time.Time `json:"first_aired" gorethink:"-" sql:"-"`
	IDs                   struct {
		Imdb   string `json:"imdb" gorethink:"-" sql:"-"`
		Tmdb   int    `json:"tmdb" gorethink:"-" sql:"-"`
		Trakt  int    `json:"trakt" gorethink:"-" sql:"-"`
		Tvdb   int    `json:"tvdb" gorethink:"-" sql:"-"`
		Tvrage int    `json:"tvrage" gorethink:"-" sql:"-"`
	} `json:"ids" gorethink:"-" sql:"-"`
	Images struct {
		Screenshot struct {
			Full   string `json:"full" gorethink:"-" sql:"-"`
			Medium string `json:"medium" gorethink:"-" sql:"-"`
			Thumb  string `json:"thumb" gorethink:"-" sql:"-"`
		} `json:"screenshot" gorethink:"-" sql:"-"`
	} `json:"images" gorethink:"-" sql:"-"`
	Number    int     `json:"number" gorethink:"number,omitempty"`
	Overview  string  `json:"overview" gorethink:"-" sql:"-"`
	Rating    float64 `json:"-" gorethink:"-" sql:"-"`
	Season    int     `json:"season" gorethink:"season,omitempty"`
	Title     string  `json:"title" gorethink:"title,omitempty"`
	UpdatedAt string  `json:"updated_at" gorethink:"-" `
	Votes     int     `json:"-" gorethink:"-" sql:"-"`

	CID   []interface{} `json:"-" gorethink:"id" sql:"-"`
	Files []*UserFile   `json:"files" gorethink:"files" sql:"-"`
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
