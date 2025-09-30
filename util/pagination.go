package util

import "errors"

// ParsePagination read args[0] dan args[1] as int
// default: limit=10, offset=0
func ParsePagination(args ...interface{}) (int, int, error) {
	limit := 10
	offset := 0

	if len(args) >= 2 {
		l, ok := args[0].(int)
		if !ok {
			return 0, 0, errors.New("limit argument is not an int")
		}
		o, ok := args[1].(int)
		if !ok {
			return 0, 0, errors.New("offset argument is not an int")
		}
		limit, offset = l, o
	}

	return limit, offset, nil
}
