package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomString(maxLength int) string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	characters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, maxLength)
	for i := range result {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}
