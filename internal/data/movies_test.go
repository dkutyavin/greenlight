package data

import (
	"maps"
	"strings"
	"testing"
	"time"

	"greenlight.dekutyavin.net/internal/validator"
)

func TestValidateMovie(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		movie      Movie
		wantErrors map[string]string
	}{
		{
			name: "valid movie",
			movie: Movie{
				Title:   "Titanic",
				Year:    1997,
				Runtime: 194,
				Genres:  []string{"melodrama", "thriller", "drama"},
			},
		},
		{
			name: "empty title",
			movie: Movie{
				Title:   "",
				Year:    1997,
				Runtime: 194,
				Genres:  []string{"melodrama"},
			},
			wantErrors: map[string]string{"title": "must be provided"},
		},
		{
			name: "title longer than 500 bytes",
			movie: Movie{
				Title:   strings.Repeat("a", 501),
				Year:    1997,
				Runtime: 194,
				Genres:  []string{"melodrama"},
			},
			wantErrors: map[string]string{"title": "must not be more than 500 bytes long"},
		},
		{
			name: "empty year",
			movie: Movie{
				Title:   "Titanic",
				Year:    0,
				Runtime: 194,
				Genres:  []string{"melodrama"},
			},
			wantErrors: map[string]string{"year": "must be provided"},
		},
		{
			name: "year too early",
			movie: Movie{
				Title:   "Titanic",
				Year:    1887,
				Runtime: 194,
				Genres:  []string{"melodrama"},
			},
			wantErrors: map[string]string{"year": "must be greater than 1888"},
		},
		{
			name: "year in the future",
			movie: Movie{
				Title:   "Titanic",
				Year:    time.Now().Year() + 1,
				Runtime: 194,
				Genres:  []string{"melodrama"},
			},
			wantErrors: map[string]string{"year": "must not be in the future"},
		},
		{
			name: "empty runtime",
			movie: Movie{
				Title:   "Titanic",
				Year:    1997,
				Runtime: 0,
				Genres:  []string{"melodrama"},
			},
			wantErrors: map[string]string{"runtime": "must be provided"},
		},
		{
			name: "negative runtime",
			movie: Movie{
				Title:   "Titanic",
				Year:    1997,
				Runtime: -194,
				Genres:  []string{"melodrama"},
			},
			wantErrors: map[string]string{"runtime": "must be positive integer"},
		},
		{
			name: "missing genres",
			movie: Movie{
				Title:   "Titanic",
				Year:    1997,
				Runtime: 194,
				Genres:  nil,
			},
			wantErrors: map[string]string{"genres": "must be provided"},
		},
		{
			name: "empty genres",
			movie: Movie{
				Title:   "Titanic",
				Year:    1997,
				Runtime: 194,
				Genres:  []string{},
			},
			wantErrors: map[string]string{"genres": "must contain at least 1 genre"},
		},
		{
			name: "more than 5 genres",
			movie: Movie{
				Title:   "Titanic",
				Year:    1997,
				Runtime: 194,
				Genres:  []string{"melodrama", "thriller", "drama", "romance", "war", "comedy"},
			},
			wantErrors: map[string]string{"genres": "must not contain more than 5 genres"},
		},
		{
			name: "duplicate genres",
			movie: Movie{
				Title:   "Titanic",
				Year:    1997,
				Runtime: 194,
				Genres:  []string{"drama", "drama"},
			},
			wantErrors: map[string]string{"genres": "must not contain duplicate genres"},
		},
		{
			name: "several invalid fields",
			movie: Movie{
				Title:   "",
				Year:    0,
				Runtime: 0,
				Genres:  []string{},
			},
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

			v := validator.New()

			ValidateMovie(v, &tt.movie)

			if tt.wantErrors == nil {
				if !v.Valid() {
					t.Errorf("want valid movie; got errors: %v", v.Errors)
				}

				return
			}

			if !maps.Equal(v.Errors, tt.wantErrors) {
				t.Errorf("want errors: %v; got errors: %v", tt.wantErrors, v.Errors)
			}
		})
	}
}
