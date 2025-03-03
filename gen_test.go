package main

import (
	"testing"

	"github.com/hex4coder/go-url-shortener/utils"
)

func TestGen(t *testing.T) {

	t.Fatal(utils.GenerateRandomString(3))
}
