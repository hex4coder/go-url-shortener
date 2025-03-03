package main

import (
	"fmt"

	"github.com/hex4coder/go-url-shortener/utils"
)

type Server struct {
	shortener utils.ShortenerI
	port      int
}

func (s *Server) Init() {

	// read file excels
	err, datalinks := s.shortener.ReadFile("./datasource/data.xlsx")

	if err != nil {
		fmt.Printf("[ERROR] - %v\r\n", err)
		return
	}

	fmt.Printf("[INFO] - Ditemukan %d data, lanjut memproses shortlinks....", len(datalinks))
	for _, dl := range datalinks {

		fmt.Println(dl.Teacher)
		fmt.Println(dl.Lesson)
		fmt.Println(dl.LongUrl)
		fmt.Println(dl.Token)
		fmt.Println("--------------------------------------------------")
		fmt.Println("")
	}
}

func (s *Server) Run() {}

func NewServer(port int, shortener utils.ShortenerI) *Server {
	return &Server{
		shortener: shortener,
		port:      port,
	}
}
