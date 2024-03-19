package utils_test

import (
	"testing"

	"github.com/ChenSongJian/ginstagram/utils"
)

func TestIsComplex(t *testing.T) {
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

func TestGenerateHash(t *testing.T) {
	password := "Password123"
	hash := utils.GenerateHash(password)
	if hash == "" {
		t.Errorf("GenerateHash(%s) returned empty string", password)
	}

	samePassword := "Password123"
	diffHash := utils.GenerateHash(samePassword)
	if hash == diffHash {
		t.Errorf("GenerateHsh should returned different hash for same password")
	}

	diffPassword := "Password456"
	diffHash = utils.GenerateHash(diffPassword)
	if hash == diffHash {
		t.Errorf("GenerateHash should returned different hash for different password")
	}
}

func TestCompareHash(t *testing.T) {
	password := "Password123"
	hash := utils.GenerateHash(password)
	result := utils.CompareHash(hash, password)
	if !result {
		t.Errorf("CompareHash(hash, password) returned false")
	}

	diffPassword := "Password456"
	result = utils.CompareHash(hash, diffPassword)
	if result {
		t.Errorf("CompareHash(hash, diffPassword) returned true")
	}
}
