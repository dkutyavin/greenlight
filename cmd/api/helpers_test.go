package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	cases := []struct {
		name   string
		data   map[string]any
		want   string
		status int
	}{
		{
			name: "valid json",
			data: map[string]any{
				"field1": "value",
				"field2": 200,
			},
			want:   `{"field1":"value","field2":200}`,
			status: http.StatusOK,
		},
	}

	wantContentType := "application/json"
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.want = tt.want + "\n"

			rr := httptest.NewRecorder()

			err := writeJSON(rr, tt.data, http.StatusOK)
			if err != nil {
				t.Fatalf("unexpected writeJSON error: %v", err)
			}

			gotContentType := rr.Header().Get("Content-Type")

			if gotContentType != wantContentType {
				t.Errorf("want content type: %q; got: %q", wantContentType, gotContentType)
			}

			if rr.Body.String() != tt.want {
				t.Errorf("want json: %q; got: %q", tt.want, rr.Body.String())
			}

		})
	}

}
