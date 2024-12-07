package models

import "errors"

type InputSong struct {
	GroupName string `json:"group" db:"group_name" binding:"required"`
	SongName  string `json:"song" db:"song_name" binding:"required"`
}

type InputUpdateSong struct {
	GroupName   *string `json:"group" db:"group_name"`
	SongName    *string `json:"song" db:"song_name"`
	ReleaseDate *string `json:"release_date" db:"release_date"`
	Text        *string `json:"text" db:"text"`
	Link        *string `json:"link" db:"link"`
}

func (i InputUpdateSong) Validate() error {
	if i.GroupName == nil && i.SongName == nil && i.ReleaseDate == nil && i.Text == nil && i.Link == nil {
		return errors.New("update structure has not values")
	}

	return nil
}
