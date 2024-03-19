package utils_test

import (
	"testing"

	"github.com/ChenSongJian/ginstagram/utils"
)

func TestIsComplex(t *testing.T) {
	// Test cases for IsComplex function
	testCases := []struct {
		password string
		expected bool
	}{
		{"Password123", true},              // Meets all requirements
		{"password123", false},             // Valid length but no uppercase letter
		{"PASSWORD123", false},             // Valid length but no lowercase letter
		{"PASSWORD", false},                // Valid length but no digit
		{"Pass123", false},                 // Valid complexity but too short
		{"ComplexAndLong123456789", false}, // Valid complexity but too long

	}

	for _, tc := range testCases {
		result := utils.IsComplex(tc.password)
		if result != tc.expected {
			t.Errorf("IsComplex(%s) = %t; expected %t", tc.password, result, tc.expected)
		}
	}
}
