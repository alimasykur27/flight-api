package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for RemoveDuplicate function
func TestRemoveDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []interface{}
	}{
		{
			"Test integer slice",
			[]interface{}{1, 2, 2, 3, 4, 4, 5},
			[]interface{}{1, 2, 3, 4, 5},
		},
		{
			"Test string slice",
			[]interface{}{"a", "b", "a", "c", "b"},
			[]interface{}{"a", "b", "c"},
		},
		{
			"Test mixed slice",
			[]interface{}{1, "a", 2, "b", 1, "a"},
			[]interface{}{1, "a", 2, "b"},
		},
		{
			"Test empty slice",
			[]interface{}{},
			[]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveDuplicate(tt.input)
			if len(result) != len(tt.expected) {
				assert.Fail(t, "Length mismatch", "Expected length %d, got %d", len(tt.expected), len(result))
			}
			assert.ElementsMatch(t, tt.expected, result, "Expected %v, got %v", tt.expected, result)

			// uniqueness check
			uniqueMap := make(map[interface{}]struct{})
			for _, item := range result {
				if _, exists := uniqueMap[item]; exists {
					assert.Fail(t, "Duplicate found in result", "Item %v is duplicated", item)
				}
				uniqueMap[item] = struct{}{}
			}
		})
	}
}
