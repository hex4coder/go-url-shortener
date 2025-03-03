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
	// new generator url
	st := utils.NewShortener("h.ei", 3)

	// read file excels
	err, _ := st.ReadFile("excelfile")

	if err != nil {
		fmt.Printf("[ERROR] - %v\r\n", err)
		return
	}

	fmt.Printf("st struc : %v \r\n", st)
}
