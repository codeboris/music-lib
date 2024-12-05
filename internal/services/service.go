package services

import (
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
