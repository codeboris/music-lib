package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	var songList = []models.Song{}

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

func (r *Repository) InsertSong(song models.Song) (int, error) {
	var id int
	createQuery := fmt.Sprintf(`
		INSERT INTO %s (group_name, song_name, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, songsTable)

	row := r.db.QueryRow(
		createQuery,
		song.GroupName,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetExternalData(url string) (models.SongDetail, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.SongDetail{}, err
	}
	defer resp.Body.Close()

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return models.SongDetail{}, err
	}

	return songDetail, nil
}

func (r *Repository) Update(songID int, input models.InputUpdateSong) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.GroupName != nil {
		setValues = append(setValues, fmt.Sprintf("group_name=$%d", argID))
		args = append(args, *&input.GroupName)
		argID++
	}

	if input.SongName != nil {
		setValues = append(setValues, fmt.Sprintf("song_name=$%d", argID))
		args = append(args, *&input.SongName)
		argID++
	}

	if input.ReleaseDate != nil {
		setValues = append(setValues, fmt.Sprintf("release_date=$%d", argID))
		args = append(args, *&input.ReleaseDate)
		argID++
	}

	if input.Text != nil {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argID))
		args = append(args, *&input.Text)
		argID++
	}

	if input.Link != nil {
		setValues = append(setValues, fmt.Sprintf("link=$%d", argID))
		args = append(args, *&input.Link)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", songsTable, setQuery, songID)
	args = append(args, songID)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *Repository) Delete(songID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", songsTable)

	_, err := r.db.Exec(query, songID)
	return err
}
