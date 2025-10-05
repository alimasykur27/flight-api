package util_test

import (
	"flight-api/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtrConverter(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "int",
			input:    42,
			expected: 42,
		},
		{
			name:     "bool",
			input:    true,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.expected.(type) {
			case string:
				result := util.Ptr(tt.input.(string))
				assert.Equal(t, v, *result, "Ptr(%v) = %v, want %v", tt.input, *result, v)
				if *result != v {
					t.Errorf("Ptr(%v) = %v, want %v", tt.input, *result, v)
				}
			case int:
				result := util.Ptr(tt.input.(int))
				assert.Equal(t, v, *result, "Ptr(%v) = %v, want %v", tt.input, *result, v)
				if *result != v {
					t.Errorf("Ptr(%v) = %v, want %v", tt.input, *result, v)
				}
			case bool:
				result := util.Ptr(tt.input.(bool))
				assert.Equal(t, v, *result, "Ptr(%v) = %v, want %v", tt.input, *result, v)
				if *result != v {
					t.Errorf("Ptr(%v) = %v, want %v", tt.input, *result, v)
				}
			}
		})
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *int64
	}{
		{
			name:     "valid input 1",
			input:    "42",
			expected: util.Ptr(int64(42)),
		},
		{
			name:     "valid input 2",
			input:    "0",
			expected: util.Ptr(int64(0)),
		},
		{
			name:     "valid input 3",
			input:    "123123123123",
			expected: util.Ptr(int64(123123123123)),
		},
		{
			name:     "valid input 4",
			input:    "-123123123123",
			expected: util.Ptr(int64(-123123123123)),
		},
		{
			name:     "invalid input 1",
			input:    "invalid",
			expected: nil,
		},
		{
			name:     "invalid input 2",
			input:    "",
			expected: nil,
		},
		{
			name:     "invalid input 3",
			input:    "123.123123123",
			expected: nil,
		},
		{
			name:     "invalid input 4",
			input:    "-123.123123123",
			expected: nil,
		},
		{
			name:     "invalid input 5",
			input:    "Ini string",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ParseInt64Ptr(tt.input)
			assert.IsType(t, tt.expected, result)
			assert.Equal(t, tt.expected, result, "ParseInt64Ptr(%v) = %v, want %v", tt.input, result, tt.expected)
			if result != nil {
				if *result != *tt.expected {
					t.Errorf("ParseInt64Ptr(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			} else {
				if tt.expected != nil {
					t.Errorf("ParseInt64Ptr(%v) = nil, want %v", tt.input, *tt.expected)
				}
			}
		})
	}
}

func TestToInterfaces(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []interface{}
	}{
		{
			name:     "valid input 1",
			input:    []string{"hello", "world"},
			expected: []interface{}{"hello", "world"},
		},
		{
			name:     "valid input 2",
			input:    []string{},
			expected: []interface{}{},
		},
		{
			name:     "valid input 3",
			input:    []string{"hello", "world", "!"},
			expected: []interface{}{"hello", "world", "!"},
		},
		{
			name:     "valid input 4",
			input:    []string{"1", "2", "3"},
			expected: []interface{}{"1", "2", "3"},
		},
		{
			name:     "valid input 5",
			input:    []string{"", "", ""},
			expected: []interface{}{"", "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ToInterfaces(tt.input)
			assert.IsType(t, tt.expected, result)
			assert.Equal(t, tt.expected, result, "ToInterfaces(%v) = %v, want %v", tt.input, result, tt.expected)
			for i, v := range result {
				assert.Equal(t, tt.expected[i], v, "ToInterfaces(%v)[%d] = %v, want %v", tt.input, i, v, tt.expected[i])
				if v != tt.expected[i] {
					t.Errorf("ToInterfaces(%v)[%d] = %v, want %v", tt.input, i,
						v, tt.expected[i])
				}
			}
		})
	}
}
