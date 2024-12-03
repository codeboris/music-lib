package services

import (
	"github.com/codeboris/music-lib/internal/models"
	"github.com/codeboris/music-lib/internal/repositories"
)

type Service struct {
	repo *repositories.SongRepository
}

func NewService(repo *repositories.SongRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FetchSongs() ([]models.Song, error) {
	return s.repo.GetSongs()
}
