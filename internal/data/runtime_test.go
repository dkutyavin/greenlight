package data

import (
	"encoding/json"
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
