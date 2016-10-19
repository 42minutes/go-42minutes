package minutes

// UserShow struct
type UserShow struct {
	ID    string `json:"id" gorethink:"id"`
	Title string `json:"title" gorethink:"title,omitempty"`
}

// UserSeason struct
type UserSeason struct {
	ShowID string `json:"show_id" gorethink:"show_id"`
	Number int    `json:"number" gorethink:"number,omitempty"`
}

// UserEpisode struct
type UserEpisode struct {
	ShowID string `json:"show_id" gorethink:"show_id"`
	Number int    `json:"number" gorethink:"number,omitempty"`
	Season int    `json:"season" gorethink:"season,omitempty"`
	Title  string `json:"title" gorethink:"title,omitempty"`
}
