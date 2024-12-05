package repositories

import (
	"fmt"
	"strings"

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

func (r *Repository) FetchSongs(filter models.SongFilter) ([]models.Song, error) {
	var songList []models.Song

	query := fmt.Sprintf("SELECT * FROM %s WHERE 1=1", songsTable)

	// Filters
	var conditions []string
	var args []interface{}
	argIdx := 1 // Index ($1, $2, ...)

	if filter.Group != "" {
		conditions = append(conditions, "group_name ILIKE $%d")
		args = append(args, "%"+filter.Group+"%")
		argIdx++
	}
	if filter.Song != "" {
		conditions = append(conditions, "song_name ILIKE $%d")
		args = append(args, "%"+filter.Song+"%")
		argIdx++
	}
	if filter.ReleaseDate != "" {
		conditions = append(conditions, "release_date = $%d")
		args = append(args, filter.ReleaseDate)
		argIdx++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Pagination
	limit := filter.Limit

	if filter.Page < 1 {
		filter.Page = 1
	}

	if limit <= 0 {
		limit = 10
	}
	offset := (filter.Page - 1) * limit
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argIdx, argIdx+1)

	args = append(args, limit, offset)

	err := r.db.Select(&songList, query, args...)
	if err != nil {
		return nil, err
	}

	return songList, nil
}
