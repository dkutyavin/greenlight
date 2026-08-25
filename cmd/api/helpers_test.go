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
