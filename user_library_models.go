package minutes

import "fmt"

// UserShow struct
type UserShow struct {
	ID    string `json:"id" gorethink:"id"`
	Title string `json:"title" gorethink:"title,omitempty"`
}

// GetCID -
func (e *UserShow) GetCID() interface{} {
	return e.ID
}

// UserSeason struct
type UserSeason struct {
	CID    []interface{} `gorethink:"id"`
	ShowID string        `json:"show_id" gorethink:"show_id"`
	Number int           `json:"number" gorethink:"number"`
}

// GetCID -
func (e *UserSeason) GetCID() []interface{} {
	return []interface{}{
		e.ShowID,
		e.Number,
	}
}

// UserEpisode struct
type UserEpisode struct {
	CID           []interface{} `gorethink:"id"`
	ShowID        string        `json:"show_id" gorethink:"show_id"`
	Season        int           `json:"season" gorethink:"season"`
	Number        int           `json:"number" gorethink:"number"`
	Title         string        `json:"title" gorethink:"title,omitempty"`
	Infohash      string        `gorethink:"infohash"`
	Downloaded    bool          `gorethink:"downloaded"`
	RetryDatetime int64         `gorethink:"retry_time"`
}

// String -
func (e *UserEpisode) String() string {
	return fmt.Sprintf("ID:%s S%02dE%02d [%s]:%t",
		e.ShowID, e.Season, e.Number, e.Infohash, e.Downloaded)
}

// GetCID -
func (e *UserEpisode) GetCID() []interface{} {
	return []interface{}{
		e.ShowID,
		e.Season,
		e.Number,
	}
}
