package models

type SongFilter struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}
