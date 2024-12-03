package repositories

import (
	"database/sql"

	"github.com/codeboris/music-lib/internal/models"
)

type SongRepository struct {
	db *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{db: db}
}

func (r *SongRepository) GetSongs() ([]models.Song, error) {
	var songs []models.Song
	songs = append(songs,
		models.Song{GroupName: "Muse", SongName: "Supermassive Black Hole"},
		models.Song{GroupName: "Muse2", SongName: "2222 Supermassive Black Hole"},
		models.Song{GroupName: "Muse3", SongName: "33333 Supermassive Black Hole"},
	)
	return songs, nil
}
