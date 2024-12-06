package services

import (
	"fmt"

	"github.com/codeboris/music-lib/internal/models"
	"github.com/codeboris/music-lib/internal/repositories"
)

type Service struct {
	repo *repositories.Repository
}

func NewService(repo *repositories.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSongList(filter models.SongFilter) ([]models.Song, error) {
	return s.repo.FetchSongs(filter)
}

func (s *Service) GetExternalData(group, song string) (models.SongDetail, error) {
	url := fmt.Sprintf("http://example.com/api/info?group=%s&song=%s", group, song)
	return s.repo.GetExternalData(url)
}

func (s *Service) CreateSong(song models.Song) (int, error) {
	return s.repo.InsertSong(song)
}

func (s *Service) PrepareSong(detail models.SongDetail, input models.InputSong) models.Song {
	return models.Song{
		GroupName:   input.GroupName,
		SongName:    input.SongName,
		ReleaseDate: detail.ReleaseDate,
		Text:        detail.Text,
		Link:        detail.Link,
	}
}

func (s *Service) UpdateSong(songID int, input models.InputUpdateSong) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(songID, input)
}

func (s *Service) DeleteSong(songID int) error {
	return s.repo.Delete(songID)
}
