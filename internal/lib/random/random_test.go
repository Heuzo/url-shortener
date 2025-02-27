package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		expected int
	}{
		{name: "length 10", size: 10, expected: 10},
		{name: "length 100", size: 100, expected: 100},
		{name: "length 1", size: 1, expected: 1},
		{name: "length 0", size: 0, expected: 0},
		{name: "length -1", size: -1, expected: 0},
		{name: "length -100", size: -100, expected: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomString(tt.size)
			str2 := NewRandomString(tt.size)

			assert.Len(t, str1, tt.expected)
			assert.Len(t, str2, tt.expected)

			if tt.expected != 0 {
				assert.NotEqual(t, str1, str2)
			}
		})
	}
}
