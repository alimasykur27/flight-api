package enum_test

import (
	"flight-api/internal/enum"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToOwnership(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected enum.OwnershipEnum
	}{
		{
			name:     "public",
			input:    "public",
			expected: enum.OWN_PUBLIC,
		},
		{
			name:     "private",
			input:    "private",
			expected: enum.OWN_PRIVATE,
		},
		{
			name:     "invalid - 1",
			input:    "invalid",
			expected: enum.OWN_NIL,
		},
		{
			name:     "invalid - 2",
			input:    "1231279y3kh",
			expected: enum.OWN_NIL,
		},
		{
			name:     "empty",
			input:    "",
			expected: enum.OWN_NIL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := enum.ToOwnership(tt.input)
			assert.IsType(t, tt.expected, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}
