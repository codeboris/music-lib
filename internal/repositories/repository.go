package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/codeboris/music-lib/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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
		conditions = append(conditions, fmt.Sprintf("group_name ILIKE $%d", argIdx))
		args = append(args, "%"+filter.Group+"%")
		argIdx++
	}
	if filter.Song != "" {
		conditions = append(conditions, fmt.Sprintf("song_name ILIKE $%d", argIdx))
		args = append(args, "%"+filter.Song+"%")
		argIdx++
	}
	if filter.ReleaseDate != "" {
		conditions = append(conditions, fmt.Sprintf("release_date = $%d", argIdx))
		args = append(args, filter.ReleaseDate)
		argIdx++
	}
	if filter.Text != "" {
		conditions = append(conditions, fmt.Sprintf("text ILIKE $%d", argIdx))
		args = append(args, "%"+filter.Text+"%")
		argIdx++
	}
	if filter.Link != "" {
		conditions = append(conditions, fmt.Sprintf("link ILIKE $%d", argIdx))
		args = append(args, "%"+filter.Link+"%")
		argIdx++
	}

	if len(conditions) > 0 {
		query += " AND (" + strings.Join(conditions, " OR ") + ")"
	}

	// Pagination
	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}

	offset := (filter.Page - 1) * limit
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argIdx, argIdx+1)

	args = append(args, limit, offset)

	logrus.Infoln(query)

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
	argIdx := 1

	if input.GroupName != nil {
		setValues = append(setValues, fmt.Sprintf("group_name=$%d", argIdx))
		args = append(args, &input.GroupName)
		argIdx++
	}

	if input.SongName != nil {
		setValues = append(setValues, fmt.Sprintf("song_name=$%d", argIdx))
		args = append(args, &input.SongName)
		argIdx++
	}

	if input.ReleaseDate != nil {
		setValues = append(setValues, fmt.Sprintf("release_date=$%d", argIdx))
		args = append(args, &input.ReleaseDate)
		argIdx++
	}

	if input.Text != nil {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argIdx))
		args = append(args, &input.Text)
		argIdx++
	}

	if input.Link != nil {
		setValues = append(setValues, fmt.Sprintf("link=$%d", argIdx))
		args = append(args, &input.Link)
		argIdx++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", songsTable, setQuery, argIdx)
	args = append(args, songID)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *Repository) Delete(songID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", songsTable)

	_, err := r.db.Exec(query, songID)
	return err
}

func (r *Repository) GetSongByID(songID int) (models.Song, error) {
	var song models.Song
	query := fmt.Sprintf("SELECT * FROM %s WHERE id= $1", songsTable)
	err := r.db.Get(&song, query, songID)
	return song, err
}
