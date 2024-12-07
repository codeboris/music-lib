package db

import (
	"context"
	"fmt"
	"time"

	"github.com/codeboris/music-lib/config"
	"github.com/jmoiron/sqlx"
)

func buildDSN(cfg config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPass, cfg.SSLMode)
}

func NewPostgresDB(cfg config.Config) (*sqlx.DB, error) {
	dsn := buildDSN(cfg)
	dbApi, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbApi.PingContext(ctx); err != nil {
		dbApi.Close()
		return nil, err
	}

	return dbApi, nil
}

// Вставка песни

// DROP TABLE IF EXISTS songs;
// var songID int
// err = db.QueryRow("INSERT INTO songs (group_name, song_name, release_date, link) VALUES ($1, $2, $3, $4) RETURNING id",
// 	"Artist Name", "Song Title", "2024-01-01", "http://example.com").Scan(&songID)
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Printf("Song added with ID: %d\n", songID)

// // Текст песни
// songText := `Ooh baby, don't you know I suffer?
// Ooh baby, can you hear me moan?
// You caught me under false pretenses
// How long before you let me go?

// Ooh
// You set my soul alight

// Ooh
// You set my soul alight`

// // Разбиваем текст песни на куплеты
// verses := strings.Split(songText, "\n\n")

// // Вставка куплетов в таблицу song_verses
// for i, verse := range verses {
// 	// Вставляем куплет в базу данных
// 	_, err := db.Exec("INSERT INTO song_verses (song_id, verse_order, verse_text) VALUES ($1, $2, $3)",
// 		songID, i+1, verse)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Verse %d inserted successfully\n", i+1)
// }

// insert song
// -- Вставка песни
// INSERT INTO songs (group_name, song_name, release_date, link)
// VALUES ('Artist Name', 'Song Title', '2024-01-01', 'http://example.com');

// insert verses
// -- Вставка куплетов
// INSERT INTO song_verses (song_id, verse_order, verse_text)
// VALUES
// (1, 1, 'Ooh baby, don\'t you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?'),
// (1, 2, 'Ooh\nYou set my soul alight'),
// (1, 3, 'Ooh\nYou set my soul alight');

// song_verses
// CREATE TABLE song_verses (
//     id SERIAL PRIMARY KEY,
//     song_id INT NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
//     verse_order INT NOT NULL,  -- Порядковый номер куплета
//     verse_text TEXT NOT NULL,  -- Текст куплета
//     CONSTRAINT unique_verse UNIQUE (song_id, verse_order)
// );

// songs
// CREATE TABLE songs (
//     id SERIAL PRIMARY KEY,
//     group_name VARCHAR(255) NOT NULL,
//     song_name VARCHAR(255) NOT NULL,
//     release_date DATE NOT NULL,
//     link VARCHAR(255),
//     CONSTRAINT song_group_unique UNIQUE (group_name, song_name)
// );
