package utils

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func IsComplex(password string) bool {
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	isValidLength := len(password) >= 8 && len(password) <= 20
	return hasLower && hasUpper && hasDigit && isValidLength
}

func GenerateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func CompareHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
