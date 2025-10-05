package util

// ParsePagination read limit and offset
// default: limit=10, offset=0
func ParsePagination(args map[string]interface{}) (int, int) {
	limit := 10
	offset := 0

	if len(args) > 0 {
		if val, ok := args["limit"]; ok {
			if v, ok := val.(int); ok {
				limit = v
			}
			if val, ok := args["offset"]; ok {
				if v, ok := val.(int); ok {
					offset = v
				}
			}
		}
	}

	return limit, offset
}
