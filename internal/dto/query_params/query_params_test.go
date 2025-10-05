package queryparams_test

import (
	queryparams "flight-api/internal/dto/query_params"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ExpextedTest struct {
	limit  int
	page   int
	offset int
}

func TestGetQueryParams(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want ExpextedTest
	}{
		{
			name: "no params -> defaults",
			url:  "/airports",
			want: ExpextedTest{limit: 10, page: 1, offset: 0},
		},
		{
			name: "valid limit and page",
			url:  "/airports?limit=50&page=3",
			want: ExpextedTest{limit: 50, page: 3, offset: 100}, // (3-1)*50
		},
		{
			name: "only limit valid",
			url:  "/airports?limit=25",
			want: ExpextedTest{limit: 25, page: 1, offset: 0},
		},
		{
			name: "only page valid",
			url:  "/airports?page=4",
			want: ExpextedTest{limit: 10, page: 4, offset: 30}, // (4-1)*10
		},
		{
			name: "non-numeric values -> defaults kept",
			url:  "/airports?limit=abc&page=xyz",
			want: ExpextedTest{limit: 10, page: 1, offset: 0},
		},
		{
			name: "zero and negative -> ignored, keep defaults",
			url:  "/airports?limit=0&page=-2",
			want: ExpextedTest{limit: 10, page: 1, offset: 0},
		},
		{
			name: "negative limit only -> ignore limit but accept valid page",
			url:  "/airports?limit=-5&page=2",
			want: ExpextedTest{limit: 10, page: 2, offset: 10}, // (2-1)*10
		},
		{
			name: "zero page only -> ignore page but accept valid limit",
			url:  "/airports?limit=7&page=0",
			want: ExpextedTest{limit: 7, page: 1, offset: 0},
		},
		{
			name: "extra unrelated params don't matter",
			url:  "/airports?limit=15&page=2&sort=desc&foo=bar",
			want: ExpextedTest{limit: 15, page: 2, offset: 15},
		},
		{
			name: "boundary small values",
			url:  "/airports?limit=1&page=1",
			want: ExpextedTest{limit: 1, page: 1, offset: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			got := queryparams.GetQueryParams(req)

			assert.Equal(t, tt.want.limit, got.Limit, "limit mismatch")
			assert.Equal(t, tt.want.page, got.Page, "page mismatch")
			assert.Equal(t, tt.want.offset, got.Offset, "offset mismatch")
		})
	}
}
