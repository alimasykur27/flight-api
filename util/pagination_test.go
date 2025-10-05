package util_test

import (
	"flight-api/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		expected struct {
			limit  int
			offset int
		}
	}{
		{
			name: "valid input 1",
			args: map[string]interface{}{
				"limit":  20,
				"offset": 3,
			},
			expected: struct {
				limit  int
				offset int
			}{
				limit:  20,
				offset: 3,
			},
		},
		{
			name: "valid input 2 - default limit",
			args: map[string]interface{}{
				"offset": 0,
			},
			expected: struct {
				limit  int
				offset int
			}{
				limit:  10,
				offset: 0,
			},
		},
		{
			name: "valid input 3 - default offset",
			args: map[string]interface{}{
				"limit": 10,
			},
			expected: struct {
				limit  int
				offset int
			}{
				limit:  10,
				offset: 0,
			},
		},
		{
			name: "valid input 4 - default limit and offset",
			args: map[string]interface{}{},
			expected: struct {
				limit  int
				offset int
			}{
				limit:  10,
				offset: 0,
			},
		},
		{
			name: "invalid input 1 - invalid limit",
			args: map[string]interface{}{
				"limit": "invalid",
			},
			expected: struct {
				limit  int
				offset int
			}{
				limit:  10,
				offset: 0,
			},
		},
		{
			name: "invalid input 2 - invalid offset",
			args: map[string]interface{}{
				"offset": 12.12,
			},
			expected: struct {
				limit  int
				offset int
			}{
				limit:  10,
				offset: 0,
			},
		},
		{
			name: "invalid input 3 - invalid limit and offset",
			args: map[string]interface{}{
				"limit":  "invalid",
				"offset": "invalid",
			},
			expected: struct {
				limit  int
				offset int
			}{
				limit:  10,
				offset: 0,
			},
		},
		{
			name: "invalid input 4 - nil args",
			args: nil,
			expected: struct {
				limit  int
				offset int
			}{
				limit:  10,
				offset: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, offset := util.ParsePagination(tt.args)
			assert.IsType(t, tt.expected.offset, offset, "ParsePagination(%v).offset = %v, want %v", tt.args, offset, tt.expected.offset)
			assert.Equal(t, tt.expected.offset, offset, "ParsePagination(%v).offset = %v, want %v", tt.args, offset, tt.expected.offset)
			assert.IsType(t, tt.expected.limit, limit, "ParsePagination(%v).limit = %v, want %v", tt.args, limit, tt.expected.limit)
			assert.Equal(t, tt.expected.limit, limit, "ParsePagination(%v).limit = %v, want %v", tt.args, limit, tt.expected.limit)
		})
	}

}
