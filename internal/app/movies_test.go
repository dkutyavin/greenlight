package app

import (
	"context"
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"greenlight.dekutyavin.net/internal/data"
)

func TestShowMovieHandler(t *testing.T) {
	t.Parallel()

	existing := data.Movie{
		ID:      1,
		Title:   "Casablanca",
		Year:    1942,
		Runtime: 102,
		Genres:  []string{"drama", "romance", "war"},
		Version: 1,
	}

	cases := []struct {
		name           string
		idParam        string
		wantMovie      bool
		getErr         error
		wantStatus     int
		wantModelCalls int
	}{
		{
			name:           "existing movie",
			idParam:        "1",
			wantMovie:      true,
			wantStatus:     http.StatusOK,
			wantModelCalls: 1,
		},
		{
			name:           "missing movie",
			idParam:        "1",
			getErr:         data.ErrRecordNotFound,
			wantStatus:     http.StatusNotFound,
			wantModelCalls: 1,
		},
		{
			name:           "model error",
			idParam:        "1",
			getErr:         errors.New("connection reset by peer"),
			wantStatus:     http.StatusInternalServerError,
			wantModelCalls: 1,
		},
		{
			name:           "invalid id",
			idParam:        "abc",
			wantStatus:     http.StatusNotFound,
			wantModelCalls: 0,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			modelCalls := 0

			model := newTestMovieModel(t)
			model.getFn = func(gotID int) (*data.Movie, error) {
				modelCalls++

				if gotID != existing.ID {
					t.Errorf("want model called with id: %d; got: %d", existing.ID, gotID)
				}

				if !tt.wantMovie {
					return nil, tt.getErr
				}

				// return a copy so handler mutations cannot alter the expected value
				found := existing
				found.Genres = slices.Clone(existing.Genres)

				return &found, nil
			}

			app := newTestApplication(t)
			app.Models.Movies = model

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/v1/movies/"+tt.idParam, nil)
			params := httprouter.Params{{Key: "id", Value: tt.idParam}}
			req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

			app.showMovieHandler(rr, req)

			if got := rr.Result().StatusCode; got != tt.wantStatus {
				t.Fatalf("want status: %d; got: %d; body: %s", tt.wantStatus, got, rr.Body.String())
			}

			if modelCalls != tt.wantModelCalls {
				t.Errorf("want model calls: %d; got: %d", tt.wantModelCalls, modelCalls)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resBody struct {
				Movie data.Movie `json:"movie"`
			}

			err := json.Unmarshal(rr.Body.Bytes(), &resBody)
			if err != nil {
				t.Fatalf("could not unmarshal body: %v", err)
			}

			if !reflect.DeepEqual(resBody.Movie, existing) {
				t.Errorf("want movie: %+v; got movie: %+v", existing, resBody.Movie)
			}
		})
	}
}

func TestCreateMovieHandler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		js             string
		insertFn       func(*data.Movie) error
		wantStatus     int
		wantLocation   string
		wantMovie      *data.Movie
		wantErrors     map[string]string
		wantErrorMsg   string
		wantModelCalls int
	}{
		{
			name: "valid movie",
			js:   `{"title": "Titanic", "year": 1997, "runtime": "194 mins", "genres": ["melodrama", "thriller", "drama"]}`,
			insertFn: func(m *data.Movie) error {
				m.ID = 42
				m.Version = 1

				return nil
			},
			wantStatus:   http.StatusCreated,
			wantLocation: "/v1/movies/42",
			wantMovie: &data.Movie{
				ID:      42,
				Title:   "Titanic",
				Year:    1997,
				Runtime: 194,
				Genres:  []string{"melodrama", "thriller", "drama"},
				Version: 1,
			},
			wantModelCalls: 1,
		},
		{
			name:       "several invalid fields",
			js:         `{"title": "", "year": 0, "runtime": "0 mins", "genres": []}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantErrors: map[string]string{
				"title":   "must be provided",
				"year":    "must be provided",
				"runtime": "must be provided",
				"genres":  "must contain at least 1 genre",
			},
			wantModelCalls: 0,
		},
		{
			name:           "malformed json",
			js:             `{"title": `,
			wantStatus:     http.StatusBadRequest,
			wantErrorMsg:   "body contains badly-formed JSON",
			wantModelCalls: 0,
		},
		{
			name:           "model error",
			js:             `{"title": "Titanic", "year": 1997, "runtime": "194 mins", "genres": ["drama"]}`,
			insertFn:       func(*data.Movie) error { return errors.New("connection reset by peer") },
			wantStatus:     http.StatusInternalServerError,
			wantModelCalls: 1,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			modelCalls := 0

			model := newTestMovieModel(t)
			model.insertFn = func(m *data.Movie) error {
				modelCalls++

				if tt.insertFn == nil {
					return nil
				}

				return tt.insertFn(m)
			}

			app := newTestApplication(t)
			app.Models.Movies = model

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/movies", strings.NewReader(tt.js))

			app.createMovieHandler(rr, req)

			if got := rr.Result().StatusCode; got != tt.wantStatus {
				t.Fatalf("want status: %d; got: %d; body: %s", tt.wantStatus, got, rr.Body.String())
			}

			if modelCalls != tt.wantModelCalls {
				t.Errorf("want model calls: %d; got: %d", tt.wantModelCalls, modelCalls)
			}

			if tt.wantLocation != "" {
				if got := rr.Result().Header.Get("Location"); got != tt.wantLocation {
					t.Errorf("want Location: %q; got: %q", tt.wantLocation, got)
				}
			}

			switch {
			case tt.wantMovie != nil:
				var resBody struct {
					Movie data.Movie `json:"movie"`
				}

				err := json.Unmarshal(rr.Body.Bytes(), &resBody)
				if err != nil {
					t.Fatalf("could not unmarshal body: %v", err)
				}

				if !reflect.DeepEqual(resBody.Movie, *tt.wantMovie) {
					t.Errorf("want movie: %+v; got movie: %+v", *tt.wantMovie, resBody.Movie)
				}

			case tt.wantErrors != nil:
				var resBody struct {
					Error map[string]string `json:"error"`
				}

				err := json.Unmarshal(rr.Body.Bytes(), &resBody)
				if err != nil {
					t.Fatalf("could not unmarshal body: %v", err)
				}

				if !maps.Equal(resBody.Error, tt.wantErrors) {
					t.Errorf("want errors: %v; got errors: %v", tt.wantErrors, resBody.Error)
				}

			case tt.wantErrorMsg != "":
				var resBody struct {
					Error string `json:"error"`
				}

				err := json.Unmarshal(rr.Body.Bytes(), &resBody)
				if err != nil {
					t.Fatalf("could not unmarshal body: %v", err)
				}

				if resBody.Error != tt.wantErrorMsg {
					t.Errorf("want error: %q; got error: %q", tt.wantErrorMsg, resBody.Error)
				}
			}
		})
	}
}
