package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

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
			Runtime string   `json:"runtime"`
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

	wantRuntime := fmt.Sprintf("%d mins", m.Runtime)
	if resBody.Movie.Runtime != wantRuntime {
		t.Errorf("want runtime: %q; got: %q", wantRuntime, resBody.Movie.Runtime)
	}

	if !slices.Equal(resBody.Movie.Genres, m.Genres) {
		t.Errorf("want movie genres: %v; got: %v", m.Genres, resBody.Movie.Genres)
	}

	if resBody.Movie.Version != m.Version {
		t.Errorf("want version: %d; got: %d", m.Version, resBody.Movie.Version)
	}

}

func TestCreateMovieHandler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		js             string
		wantBody       string
		wantStatusCode int
	}{
		{
			name:           "correct request json",
			js:             `{"title": "Titanic", "year": 1997, "runtime": "194 mins", "genres": ["melodrama", "thriller", "drama"]}`,
			wantBody:       "{Title:Titanic Year:1997 Runtime:194 Genres:[melodrama thriller drama]}\n",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "empty title",
			js:             `{"title": "", "year": 1997, "runtime": "194 mins", "genres": ["melodrama", "thriller", "drama"]}`,
			wantBody:       `"title": "must be provided"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "title longer than 500 bytes",
			js:             `{"title": "` + strings.Repeat("a", 501) + `", "year": 1997, "runtime": "194 mins", "genres": ["melodrama"]}`,
			wantBody:       `"title": "must not be more than 500 bytes long"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "empty year",
			js:             `{"title": "Titanic", "year": 0, "runtime": "194 mins", "genres": ["melodrama"]}`,
			wantBody:       `"year": "must be provided"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "year too early",
			js:             `{"title": "Titanic", "year": 1887, "runtime": "194 mins", "genres": ["melodrama"]}`,
			wantBody:       `"year": "must be greater than 1888"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "year in the future",
			js:             fmt.Sprintf(`{"title": "Titanic", "year": %d, "runtime": "194 mins", "genres": ["melodrama"]}`, time.Now().Year()+1),
			wantBody:       `"year": "must not be in the future"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "empty runtime",
			js:             `{"title": "Titanic", "year": 1997, "runtime": "0 mins", "genres": ["melodrama"]}`,
			wantBody:       `"runtime": "must be provided"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "negative runtime",
			js:             `{"title": "Titanic", "year": 1997, "runtime": "-194 mins", "genres": ["melodrama"]}`,
			wantBody:       `"runtime": "must be positive integer"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "missing genres",
			js:             `{"title": "Titanic", "year": 1997, "runtime": "194 mins"}`,
			wantBody:       `"genres": "must be provided"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "empty genres",
			js:             `{"title": "Titanic", "year": 1997, "runtime": "194 mins", "genres": []}`,
			wantBody:       `"genres": "must contain at least 1 genre"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "more than 5 genres",
			js:             `{"title": "Titanic", "year": 1997, "runtime": "194 mins", "genres": ["melodrama", "thriller", "drama", "romance", "war", "comedy"]}`,
			wantBody:       `"genres": "must not contain more than 5 genres"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "duplicate genres",
			js:             `{"title": "Titanic", "year": 1997, "runtime": "194 mins", "genres": ["drama", "drama"]}`,
			wantBody:       `"genres": "must not contain duplicate genres"`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/movies", strings.NewReader(tt.js))

			app := newTestApplication(t)

			app.createMovieHandler(rr, req)

			gotBody := rr.Body.String()

			if !strings.Contains(gotBody, tt.wantBody) {
				t.Errorf("want body: %s; got: %s", tt.wantBody, gotBody)
			}

			if gotCode := rr.Result().StatusCode; gotCode != tt.wantStatusCode {
				t.Errorf("want status code: %d; got status code: %d", tt.wantStatusCode, gotCode)
			}
		})
	}
}
