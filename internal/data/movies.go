package data

import "time"

type Movie struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int       `json:"year,omitzero"`
	Runtime   int       `json:"runtime,omitzero"`
	Genres    []string  `json:"genres"`
	Version   int       `json:"version"`
}
