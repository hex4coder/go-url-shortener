package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomString(maxLength int) string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	characters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	length := rand.Intn(maxLength) + 1 // Generate a random length between 1 and maxLength
	result := make([]rune, length)
	for i := range result {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}
