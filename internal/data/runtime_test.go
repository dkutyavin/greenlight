package data

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestRuntimeMarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		runtime Runtime
		want    string
	}{
		{
			name:    "case 1",
			runtime: Runtime(100),
			want:    "\"100 mins\"",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			js, err := json.Marshal(tt.runtime)
			if err != nil {
				t.Fatalf("could not marshal runtime: %v", err)
			}

			if string(js) != tt.want {
				t.Errorf("want runtime: %q; got: %q", tt.want, string(js))
			}
		})
	}

}

func TestRuntimeUnmarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		body    string
		want    Runtime
		wantErr error
	}{
		{
			name: "correct runtime body",
			body: "\"100 mins\"",
			want: Runtime(100),
		},
		{
			name:    "body with incorrect format",
			body:    "100",
			wantErr: ErrInvalidRuntimeFormat,
		},
		{
			name:    "body with three parts",
			body:    `"100 mins plus"`,
			wantErr: ErrInvalidRuntimeFormat,
		},
		{
			name:    "body with misplaced parts",
			body:    `"mins 100"`,
			wantErr: ErrInvalidRuntimeFormat,
		},
		{
			name:    "body with incorrect mins part",
			body:    `"100 minutes"`,
			wantErr: ErrInvalidRuntimeFormat,
		},
		{
			name:    "body with incorrect number part",
			body:    `"hundred minutes"`,
			wantErr: ErrInvalidRuntimeFormat,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var dst Runtime
			err := json.Unmarshal([]byte(tt.body), &dst)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("wantErr: %s; gotErr: %s", err.Error(), tt.wantErr.Error())
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %q", err.Error())
			}

			if dst != tt.want {
				t.Errorf("want: %d; got: %d", tt.want, dst)
			}
		})
	}

}
