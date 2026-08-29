package app

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.dekutyavin.net/internal/data"
)

func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new movie")
}

func (app *Application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, movie, http.StatusOK)
	if err != nil {
		app.Logger.Error(err.Error())
		http.Error(w, "The server encoutered a problem and could not process your request", http.StatusInternalServerError)
	}
}
