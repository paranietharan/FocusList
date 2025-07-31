package utils

import (
	"math/rand"
	"time"
)

func GenerateVerificationCode(length int) string {
	digits := "0123456789"
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}
