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

func TestMovieUpdateHandler(t *testing.T) {
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
		name            string
		idParam         string
		js              string
		getFn           func(int) (*data.Movie, error)
		updateFn        func(*data.Movie) error
		wantStatus      int
		wantMovie       *data.Movie
		wantErrors      map[string]string
		wantGetCalls    int
		wantUpdateCalls int
	}{
		{
			name:    "existing json updates if correct body supplied",
			idParam: "1",
			js:      `{"title": "Casablanca", "year": 1942, "runtime": "102 mins", "genres": ["drama", "romance"]}`,
			updateFn: func(m *data.Movie) error {
				// the model rewrites the record; 7 comes neither from the request nor from Get
				m.Version = 7

				return nil
			},
			wantStatus: http.StatusOK,
			wantMovie: &data.Movie{
				ID:      1,
				Title:   "Casablanca",
				Year:    1942,
				Runtime: 102,
				Genres:  []string{"drama", "romance"},
				Version: 7,
			},
			wantGetCalls:    1,
			wantUpdateCalls: 1,
		},
		{
			name:            "id param is not a number",
			idParam:         "abc",
			js:              `{}`,
			wantStatus:      http.StatusNotFound,
			wantGetCalls:    0,
			wantUpdateCalls: 0,
		},
		{
			name:            "id param is missing",
			idParam:         "",
			js:              `{}`,
			wantStatus:      http.StatusNotFound,
			wantGetCalls:    0,
			wantUpdateCalls: 0,
		},
		{
			name:            "model get error",
			idParam:         "1",
			js:              `{}`,
			getFn:           func(int) (*data.Movie, error) { return nil, errors.New("connection reset by peer") },
			wantStatus:      http.StatusInternalServerError,
			wantGetCalls:    1,
			wantUpdateCalls: 0,
		},
		{
			name:            "malformed json",
			idParam:         "1",
			js:              `{"title": `,
			wantStatus:      http.StatusBadRequest,
			wantGetCalls:    1,
			wantUpdateCalls: 0,
		},
		{
			name:       "invalid fields",
			idParam:    "1",
			js:         `{"title": "", "year": 0, "runtime": "0 mins", "genres": []}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantErrors: map[string]string{
				"title":   "must be provided",
				"year":    "must be provided",
				"runtime": "must be provided",
				"genres":  "must contain at least 1 genre",
			},
			wantGetCalls:    1,
			wantUpdateCalls: 0,
		},
		{
			name:            "model update error",
			idParam:         "1",
			js:              `{"title": "Casablanca", "year": 1942, "runtime": "102 mins", "genres": ["drama"]}`,
			updateFn:        func(*data.Movie) error { return errors.New("connection reset by peer") },
			wantStatus:      http.StatusInternalServerError,
			wantGetCalls:    1,
			wantUpdateCalls: 1,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			getCalls := 0
			updateCalls := 0

			model := newTestMovieModel(t)
			model.getFn = func(id int) (*data.Movie, error) {
				getCalls++

				if tt.getFn != nil {
					return tt.getFn(id)
				}

				if id == existing.ID {
					// return a copy so handler mutations cannot alter the expected value
					found := existing
					found.Genres = slices.Clone(existing.Genres)

					return &found, nil
				}

				return nil, data.ErrRecordNotFound
			}
			model.updateFn = func(m *data.Movie) error {
				updateCalls++

				if tt.updateFn != nil {
					return tt.updateFn(m)
				}

				return nil
			}

			app := newTestApplication(t)
			app.Models.Movies = model

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/v1/movies/"+tt.idParam, strings.NewReader(tt.js))

			if tt.idParam != "" {
				params := httprouter.Params{{Key: "id", Value: tt.idParam}}
				req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
			}

			app.updateMovieHandler(rr, req)

			if gotStatus := rr.Result().StatusCode; gotStatus != tt.wantStatus {
				t.Fatalf("want status code: %d; got: %d; body: %s", tt.wantStatus, gotStatus, rr.Body.String())
			}

			if getCalls != tt.wantGetCalls {
				t.Errorf("want model Get calls: %d; got: %d", tt.wantGetCalls, getCalls)
			}

			if updateCalls != tt.wantUpdateCalls {
				t.Errorf("want model Update calls: %d; got: %d", tt.wantUpdateCalls, updateCalls)
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
			}
		})
	}
}

func TestDeleteMovieHandler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		idParam         string
		deleteFn        func(int) error
		wantStatus      int
		wantMessage     string
		wantID          int
		wantDeleteCalls int
	}{
		{
			name:            "existing movie",
			idParam:         "1",
			wantStatus:      http.StatusOK,
			wantMessage:     "movie successfullly deleted",
			wantID:          1,
			wantDeleteCalls: 1,
		},
		{
			name:            "id param is not a number",
			idParam:         "abc",
			wantStatus:      http.StatusNotFound,
			wantDeleteCalls: 0,
		},
		{
			name:            "id param is missing",
			idParam:         "",
			wantStatus:      http.StatusNotFound,
			wantDeleteCalls: 0,
		},
		{
			name:            "movie not found",
			idParam:         "1",
			deleteFn:        func(int) error { return data.ErrRecordNotFound },
			wantStatus:      http.StatusNotFound,
			wantID:          1,
			wantDeleteCalls: 1,
		},
		{
			name:            "model error",
			idParam:         "1",
			deleteFn:        func(int) error { return errors.New("connection reset by peer") },
			wantStatus:      http.StatusInternalServerError,
			wantID:          1,
			wantDeleteCalls: 1,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deleteCalls := 0

			model := newTestMovieModel(t)
			model.deleteFn = func(id int) error {
				deleteCalls++

				if tt.wantID != 0 && id != tt.wantID {
					t.Errorf("want model called with id: %d; got: %d", tt.wantID, id)
				}

				if tt.deleteFn != nil {
					return tt.deleteFn(id)
				}

				return nil
			}

			app := newTestApplication(t)
			app.Models.Movies = model

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/v1/movies/"+tt.idParam, nil)

			if tt.idParam != "" {
				params := httprouter.Params{{Key: "id", Value: tt.idParam}}
				req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
			}

			app.deleteMovieHandler(rr, req)

			if got := rr.Result().StatusCode; got != tt.wantStatus {
				t.Fatalf("want status: %d; got: %d; body: %s", tt.wantStatus, got, rr.Body.String())
			}

			if deleteCalls != tt.wantDeleteCalls {
				t.Errorf("want model Delete calls: %d; got: %d", tt.wantDeleteCalls, deleteCalls)
			}

			if tt.wantMessage == "" {
				return
			}

			var resBody struct {
				Message string `json:"message"`
			}

			if err := json.Unmarshal(rr.Body.Bytes(), &resBody); err != nil {
				t.Fatalf("could not unmarshal body: %v", err)
			}

			if resBody.Message != tt.wantMessage {
				t.Errorf("want message: %q; got message: %q", tt.wantMessage, resBody.Message)
			}
		})
	}
}
