package repositories

import (
	"fmt"

	"github.com/codeboris/music-lib/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	songsTable = "songs"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FetchSongs() ([]models.Song, error) {
	var songList []models.Song

	query := fmt.Sprintf("SELECT * FROM %s", songsTable)
	err := r.db.Select(&songList, query)

	return songList, err
}
