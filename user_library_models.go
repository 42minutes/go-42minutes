package minutes

import "fmt"

// UserShow struct
type UserShow struct {
	ID    string `json:"id" gorethink:"id"`
	Title string `json:"title" gorethink:"title,omitempty"`
}

// UserSeason struct
type UserSeason struct {
	ShowID string `json:"show_id" gorethink:"id[0]"`
	Number int    `json:"number" gorethink:"id[1]"`
}

// UserEpisode struct
type UserEpisode struct {
	ShowID        string `json:"show_id" gorethink:"id[0]"`
	Season        int    `json:"season" gorethink:"id[1]"`
	Number        int    `json:"number" gorethink:"id[2]"`
	Title         string `json:"title" gorethink:"title,omitempty"`
	Infohash      string `gorethink:"infohash"`
	Downloaded    bool   `gorethink:"downloaded"`
	RetryDatetime int64  `gorethink:"retry_time"`
}

func (e *UserEpisode) String() string {
	return fmt.Sprintf("ID:%s S%02dE%02d [%s]:%t",
		e.ShowID, e.Season, e.Number, e.Infohash, e.Downloaded)
}
