package queryparams

import (
	"net/http"
	"strconv"
)

type QueryParams struct {
	Limit  int
	Page   int
	Offset int
}

// GetPagination baca query param `limit` & `page` dari request.
// default: limit=10, page=1
func GetQueryParams(r *http.Request) QueryParams {
	query := r.URL.Query()

	limit := 10
	page := 1

	if lStr := query.Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	if pStr := query.Get("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	}

	offset := (page - 1) * limit

	return QueryParams{
		Limit:  limit,
		Page:   page,
		Offset: offset,
	}
}
