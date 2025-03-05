package main

import (
	"fmt"
	"time"

	"github.com/hex4coder/go-url-shortener/utils"
)

func main() {
	// start program
	start := time.Now()
	defer func() {
		fmt.Println("I am done")
		end := time.Since(start)
		fmt.Printf("[INFO] - Program running in : %v\r\n", end)
	}()

	// setup database
	InitDB()

	// new generator url
	st := utils.NewShortener(DB, "h.ei", 3)

	// new server
	server := NewServer(WithPort(1994), WithShortener(st))
	server.Init()
	server.Run()
}
