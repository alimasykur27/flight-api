package enum_test

import (
	"flight-api/internal/enum"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFacilityType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected enum.FasilityTypeEnum
	}{
		{
			name:     "airport",
			input:    "airport",
			expected: enum.AIRPORT,
		},
		{
			name:     "heliport",
			input:    "heliport",
			expected: enum.HELIPORT,
		},
		{
			name:     "invalid - 1",
			input:    "invalid",
			expected: enum.NIL,
		},
		{
			name:     "invalid - 2",
			input:    "1231279y3kh",
			expected: enum.NIL,
		},
		{
			name:     "empty",
			input:    "",
			expected: enum.NIL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := enum.ToFacilityType(tt.input)
			assert.IsType(t, tt.expected, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}
