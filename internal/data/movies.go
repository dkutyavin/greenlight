package data

import "time"

type Movie struct {
	ID        int       // Unique integer ID for the movie
	CreatedAt time.Time // Timestamp for when the movie is added to our database
	Title     string    // Movie title
	Year      int       // Movie release year
	Runtime   int       // Movie runtim (in minutes)
	Genres    []string  // Slice of genres for the movie (romance, comedy, etc.)
	Version   int       // The version number starts at 1 and will be incremented each time the movie information is updated
}
