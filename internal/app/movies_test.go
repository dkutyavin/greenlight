package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"testing"

	"github.com/julienschmidt/httprouter"
	"greenlight.dekutyavin.net/internal/data"
)

func TestShowMovieHandler(t *testing.T) {
	t.Parallel()

	m := data.Movie{
		ID:      1,
		Title:   "Casablanca",
		Runtime: 102,
		Genres:  []string{"drama", "romance", "war"},
		Version: 1,
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/movies/%d", m.ID), nil)
	params := httprouter.Params{{Key: "id", Value: strconv.Itoa(m.ID)}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

	app := newTestApplication(t)

	app.showMovieHandler(rr, req)

	var resBody struct {
		Movie struct {
			ID      int      `json:"id"`
			Title   string   `json:"title"`
			Runtime int      `json:"runtime"`
			Genres  []string `json:"genres"`
			Version int      `json:"version"`
		} `json:"movie"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &resBody)
	if err != nil {
		t.Fatalf("could not unmarshal body: %v", err)
	}

	if resBody.Movie.ID != m.ID {
		t.Errorf("want movie id: %d; got: %d", m.ID, resBody.Movie.ID)
	}

	if resBody.Movie.Title != m.Title {
		t.Errorf("want movie title: %q; got: %q", m.Title, resBody.Movie.Title)
	}

	if resBody.Movie.Runtime != m.Runtime {
		t.Errorf("want runtime: %d; got: %d", m.Runtime, resBody.Movie.Runtime)
	}

	if !slices.Equal(resBody.Movie.Genres, m.Genres) {
		t.Errorf("want movie genres: %v; got: %v", m.Genres, resBody.Movie.Genres)
	}

	if resBody.Movie.Version != m.Version {
		t.Errorf("want version: %d; got: %d", m.Version, resBody.Movie.Version)
	}
}
