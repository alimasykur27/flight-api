package util

// ParsePagination read limit and offset
// default: limit=10, offset=0
func ParsePagination(args map[string]interface{}) (int, int) {
	limit := 10
	offset := 0

	if len(args) > 0 {
		if val, ok := args["limit"]; ok {
			limit = val.(int)
		}

		if val, ok := args["offset"]; ok {
			offset = val.(int)
		}
	}

	return limit, offset
}
