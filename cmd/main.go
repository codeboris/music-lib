package main

import (
	"log"

	"github.com/codeboris/music-lib/config"
	"github.com/codeboris/music-lib/internal/handlers"
	"github.com/codeboris/music-lib/internal/repositories"
	"github.com/codeboris/music-lib/internal/services"
	"github.com/codeboris/music-lib/pkg/db"
	"github.com/codeboris/music-lib/pkg/server"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// @title Music Lib App API
// @version 1.0
// @description API Server for Music Lib Application

// @host localhost:{APP_PORT}
// @BasePath /api
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	cfg := config.LoadConfig()

	dbApi, err := db.NewPostgresDB(cfg)
	if err != nil {
		logrus.Fatalf("Ошибка подключения к БД: %v", err)
	}

	if err := db.RunMigrations(dbApi.DB, cfg.ConfigMigrate); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	repo := repositories.NewRepository(dbApi)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)

	srv := new(server.Server)
	if err := srv.Run(cfg.Port, handler.InitRoutes()); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
