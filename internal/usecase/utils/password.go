package utils

import (
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const MaxPasswordBytes = 72

func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPassword(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func ValidatePassword(password, email string) string {
	switch {
	case len(password) < 8:
		return "Minimal 8 karakter"
	case len(password) > MaxPasswordBytes:
		return "Kata sandi terlalu panjang"
	case !hasLetterAndDigit(password):
		return "Harus berisi huruf dan angka"
	case strings.EqualFold(password, email):
		return "Tidak boleh sama dengan email"
	}
	return ""
}

func hasLetterAndDigit(s string) bool {
	var letter, digit bool
	for _, r := range s {
		letter = letter || unicode.IsLetter(r)
		digit = digit || unicode.IsDigit(r)
	}
	return letter && digit
}
