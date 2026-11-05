package data

import (
	"maps"
	"testing"

	"greenlight.dekutyavin.net/internal/validator"
)

func TestValidateFilters(t *testing.T) {
	t.Parallel()
	safelist := []string{"id", "title"}

	cases := []struct {
		name       string
		filters    Filters
		wantValid  bool
		wantErrors map[string]string
	}{
		{
			name: "valid filter",
			filters: Filters{
				Page:         1,
				PageSize:     10,
				Sort:         "id",
				SortSafelist: safelist,
			},
			wantValid: true,
		},
		{
			name: "page equals to zero",
			filters: Filters{
				Page:         0,
				PageSize:     10,
				Sort:         "id",
				SortSafelist: safelist,
			},
			wantValid: false,
			wantErrors: map[string]string{
				"page": "must be greater than zero",
			},
		},
		{
			name: "page greater than ten million",
			filters: Filters{
				Page:         10_000_001,
				PageSize:     10,
				Sort:         "id",
				SortSafelist: safelist,
			},
			wantValid: false,
			wantErrors: map[string]string{
				"page": "must be a maximum of 10 million",
			},
		},
		{
			name: "page at upper bound",
			filters: Filters{
				Page:         10_000_000,
				PageSize:     10,
				Sort:         "id",
				SortSafelist: safelist,
			},
			wantValid: true,
		},
		{
			name: "page_size equals to zero",
			filters: Filters{
				Page:         1,
				PageSize:     0,
				Sort:         "id",
				SortSafelist: safelist,
			},
			wantValid: false,
			wantErrors: map[string]string{
				"page_size": "must be greater than zero",
			},
		},
		{
			name: "page_size greater than one hundred",
			filters: Filters{
				Page:         1,
				PageSize:     101,
				Sort:         "id",
				SortSafelist: safelist,
			},
			wantValid: false,
			wantErrors: map[string]string{
				"page_size": "must be a maximum of 100",
			},
		},
		{
			name: "page_size at upper bound",
			filters: Filters{
				Page:         1,
				PageSize:     100,
				Sort:         "id",
				SortSafelist: safelist,
			},
			wantValid: true,
		},
		{
			name: "sort is not permitted",
			filters: Filters{
				Page:         1,
				PageSize:     10,
				Sort:         "created_at",
				SortSafelist: safelist,
			},
			wantValid: false,
			wantErrors: map[string]string{
				"sort": "invalid sort value",
			},
		},
		{
			name: "several invalid filters",
			filters: Filters{
				Page:         0,
				PageSize:     101,
				Sort:         "created_at",
				SortSafelist: safelist,
			},
			wantValid: false,
			wantErrors: map[string]string{
				"page":      "must be greater than zero",
				"page_size": "must be a maximum of 100",
				"sort":      "invalid sort value",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			v := validator.New()
			ValidateFilters(v, tc.filters)

			if v.Valid() != tc.wantValid {
				t.Errorf("valid: want %t; got %t", tc.wantValid, v.Valid())
			}

			if !maps.Equal(v.Errors, tc.wantErrors) {
				t.Errorf("errors: want %v, got %v", tc.wantErrors, v.Errors)
			}
		})
	}
}
