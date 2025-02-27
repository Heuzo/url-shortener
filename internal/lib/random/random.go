package random

import (
	"math/rand"
)

const ASCII = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func NewRandomString(length int) string {
	if length <= 0 {
		return ""
	}
	bytes := make([]byte, length)
	for idx := range bytes {
		bytes[idx] = getRandomChar()
	}
	return string(bytes)
}

func getRandomChar() byte {
	return ASCII[rand.Intn(len(ASCII))]
}
