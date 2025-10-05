package enum_test

import (
	"flight-api/internal/enum"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUseType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected enum.UseTypeEnum
	}{
		{
			name:     "public",
			input:    "public",
			expected: enum.USE_PUBLIC,
		},
		{
			name:     "private",
			input:    "private",
			expected: enum.USE_PRIVATE,
		},
		{
			name:     "invalid - 1",
			input:    "invalid",
			expected: enum.USE_NIL,
		},
		{
			name:     "invalid - 2",
			input:    "1231279y3kh",
			expected: enum.USE_NIL,
		},
		{
			name:     "empty",
			input:    "",
			expected: enum.USE_NIL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := enum.ToUseType(tt.input)
			assert.IsType(t, tt.expected, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}
