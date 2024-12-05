package models

type SongLyrics struct {
	ID          int    `json:"id" db:"id"`
	SongID      string `json:"song_id" db:"song_id"`
	VerseNumber string `json:"verse_number" db:"verse_number"`
	Lyrics      string `json:"lyrics" db:"lyrics"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
}
