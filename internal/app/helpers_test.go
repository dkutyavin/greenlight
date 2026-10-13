package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestReadIDParam(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		id      string
		want    int
		wantErr error
	}{
		{
			name: "valid param",
			id:   "100",
			want: 100,
		},
		{
			name:    "param with string",
			id:      "not_a_valid_id",
			wantErr: errInvalidIdParameter,
		},
		{
			name:    "negative id",
			id:      "-1",
			wantErr: errInvalidIdParameter,
		},
		{
			name:    "zero id",
			id:      "0",
			wantErr: errInvalidIdParameter,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/movies/%s", tt.id), nil)
			params := httprouter.Params{{Key: "id", Value: tt.id}}
			req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

			got, gotErr := readIDParam(req)

			if tt.wantErr == nil && gotErr != nil {
				t.Errorf("unexpected err: %s", gotErr.Error())
			}

			if tt.wantErr != nil && !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("wanterr: %s, goterr: %s", tt.wantErr.Error(), gotErr.Error())
			}

			if got != tt.want {
				t.Errorf("want = %d, got = %d", 100, got)
			}

		})
	}

}

func TestWriteJson(t *testing.T) {
	t.Parallel()

	app := newTestApplication(t)

	cases := []struct {
		name       string
		data       envelope
		status     int
		want       string
		wantStatus int
	}{
		{
			name: "map value",
			data: map[string]any{
				"field1": "value",
				"field2": 200,
			},
			status:     http.StatusOK,
			want:       "{\n\t\"field1\": \"value\",\n\t\"field2\": 200\n}",
			wantStatus: http.StatusOK,
		},
		{
			name:       "status 500",
			data:       map[string]any{"key": "value"},
			status:     http.StatusBadRequest,
			want:       "{\n\t\"key\": \"value\"\n}",
			wantStatus: http.StatusOK,
		},
	}

	wantContentType := "application/json"
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.want = tt.want + "\n"

			rr := httptest.NewRecorder()

			err := app.writeJSON(rr, tt.data, http.StatusOK)
			if err != nil {
				t.Fatalf("unexpected writeJSON error: %v", err)
			}

			gotStatus := rr.Result().StatusCode
			gotContentType := rr.Header().Get("Content-Type")
			gotBody := rr.Body.String()

			if gotContentType != wantContentType {
				t.Errorf("want content type: %q; got: %q", wantContentType, gotContentType)
			}

			if gotStatus != tt.wantStatus {
				t.Errorf("want status: %v; got: %v", tt.wantStatus, gotStatus)
			}

			if gotBody != tt.want {
				t.Errorf("want json: %q; got: %q", tt.want, gotBody)
			}

		})
	}

}

func TestReadJson(t *testing.T) {

	type movie struct {
		Title string `json:"title"`
		Year  int    `json:"year"`
	}

	cases := []struct {
		name       string
		body       string
		want       movie
		wantErrMsg string
	}{
		{
			name: "valid json",
			body: `{"title": "Moana", "year": 2016}`,
			want: movie{Title: "Moana", Year: 2016},
		},
		{
			name:       "json with non-wrapped keys",
			body:       `{title:"Moana"}`,
			wantErrMsg: "body contains badly-formed JSON (at character 2)",
		},
		{
			name:       "truncated json",
			body:       `{"title":"Moana"`,
			wantErrMsg: "body contains badly-formed JSON",
		},
		{
			name:       "json with stringified year",
			body:       `{"title": "Moana", "year": "1984"}`,
			wantErrMsg: `body contains incorrect JSON type for field "year"`,
		},
		{
			name:       "empty body",
			body:       "",
			wantErrMsg: "body must not be empty",
		},
		{
			name:       "body larget that 1 MB",
			body:       fmt.Sprintf(`{"title": "%s", "year": 1984}`, strings.Repeat("a", 1024*1024+1)),
			wantErrMsg: "body must not be larger than 1048576 bytes",
		},
		{
			name:       "body with unknown field",
			body:       `{"title": "Moana", "year": 1984, "genres": ["adventures"]}`,
			wantErrMsg: "body contains unknown key \"genres\"",
		},
		{
			name:       "body with two jsons",
			body:       `{"title": "Moana", "year": 2016}{"title": "Moana", "year": 2016}`,
			wantErrMsg: "body must only contain a single JSON value",
		},
	}

	app := newTestApplication(t)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var dst movie

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))

			err := app.readJSON(rr, req, &dst)

			if tc.wantErrMsg != "" {
				if err == nil {
					t.Fatalf("want error %q; got nil", tc.wantErrMsg)
				}

				if !strings.Contains(err.Error(), tc.wantErrMsg) {
					t.Errorf("want error containing %q; got %q", tc.wantErrMsg, err.Error())
				}

				return
			}

			if err != nil {
				t.Fatalf("unexptected errof: %v", err)
			}

			if dst != tc.want {
				t.Errorf("want: %v; got: %v", tc.want, dst)
			}
		})
	}

}
