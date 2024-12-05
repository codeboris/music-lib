CREATE TABLE song_lyrics (
    id SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id) ON DELETE CASCADE,
    verse_number INT NOT NULL,
    lyrics TEXT NOT NULL
);
