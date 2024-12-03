package main

import (
	"log"
	"os"

	"github.com/codeboris/music-lib/internal/handlers"
	"github.com/codeboris/music-lib/internal/repositories"
	"github.com/codeboris/music-lib/internal/services"
	"github.com/codeboris/music-lib/pkg/db"
	"github.com/codeboris/music-lib/pkg/server"
	_ "github.com/lib/pq"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("APP_PORT")
	dbConn, err := db.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	repo := repositories.NewSongRepository(dbConn)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)

	srv := new(server.Server)
	if err := srv.Run(port, handler.InitRoutes()); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
