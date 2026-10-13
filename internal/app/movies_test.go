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
		wantStatusCode int
		wantBody       string
		wantErrors     map[string]string
	}{
		{
			name:           "correct request json",
			js:             `{"title": "Titanic", "year": 1997, "runtime": "194 mins", "genres": ["melodrama", "thriller", "drama"]}`,
			wantBody:       "{Title:Titanic Year:1997 Runtime:194 Genres:[melodrama thriller drama]}\n",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "several invalid fields",
			js:             `{"title": "", "year": 0, "runtime": "0 mins", "genres": []}`,
			wantStatusCode: http.StatusUnprocessableEntity,
			wantErrors: map[string]string{
				"title":   "must be provided",
				"year":    "must be provided",
				"runtime": "must be provided",
				"genres":  "must contain at least 1 genre",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/movies", strings.NewReader(tt.js))

			app := newTestApplication(t)

			app.createMovieHandler(rr, req)

			if tt.wantErrors != nil {
				assertErrorResponse(t, rr, tt.wantStatusCode, tt.wantErrors)
				return
			}

			if gotCode := rr.Result().StatusCode; gotCode != tt.wantStatusCode {
				t.Errorf("want status code: %d; got status code: %d", tt.wantStatusCode, gotCode)
			}

			if gotBody := rr.Body.String(); gotBody != tt.wantBody {
				t.Errorf("want body: %q; got: %q", tt.wantBody, gotBody)
			}
		})
	}
}
