package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hex4coder/go-url-shortener/models"
	"github.com/hex4coder/go-url-shortener/utils"
)

type Server struct {
	releasemode bool
	shortener   utils.ShortenerI
	port        int
	shortlinks  []*models.ShortLink
	ssl         bool
	verbose     bool
	domain      string
}
type ServerResponse struct {
	s       *Server
	handler http.Handler
}

func NewServer(port int, shortener utils.ShortenerI, verbose bool, ssl bool) *Server {
	return &Server{
		shortener: shortener,
		port:      port,
		domain:    "he.i",
		ssl:       ssl,
		verbose:   verbose,
	}
}
func (s *Server) Init() {

	// read file excels
	err, datalinks := s.shortener.ReadFile("./datasource/data.xlsx")

	if err != nil {
		fmt.Printf("[ERROR] - %v\r\n", err)
		return
	}

	fmt.Printf("[INFO] - Ditemukan %d data, lanjut memproses shortlinks....", len(datalinks))
	if s.verbose {
		for _, dl := range datalinks {

			fmt.Println(dl.Teacher)
			fmt.Println(dl.Lesson)
			fmt.Println(dl.LongUrl)
			fmt.Println(dl.Token)
			fmt.Println("--------------------------------------------------")
			fmt.Println("")
		}
	}

	err, shortlinks := s.shortener.GenerateShortLink(datalinks)
	if err != nil {

		fmt.Printf("[ERROR] - %v\r\n", err)
		return

	}

	basehttp := "http"
	if s.ssl {
		basehttp = "https"
	}

	for i, link := range shortlinks {
		domain := s.domain

		if s.releasemode == false {
			domain = fmt.Sprintf("localhost:%d", s.port)
		}

		link.ShortUrl = fmt.Sprintf("%s://%s/%s", basehttp, domain, link.UniqueCode)
		shortlinks[i] = link
	}

	s.shortlinks = shortlinks
	fmt.Printf("[GENERATED] - Berhasil membuat link pendek sebanyak : %d link\r\n", len(shortlinks))
}

func rootHandler(s *Server) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: Create browser parser only for Exam Browser
		fmt.Printf("Incoming request from : r.UserAgent(): %v\n", r.UserAgent()) // Get the full URL path
		path := r.URL.Path

		// Remove the leading slash
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}

		// Split the path into segments
		segments := strings.Split(path, "/")
		if len(segments) > 0 && segments[0] != "" {
			// Access specific parameters.
			uniqId := segments[0]
			found := false
			flink := new(models.ShortLink)
			for _, link := range s.shortlinks {
				if link.UniqueCode == uniqId {
					found = true
					flink = link
					break
				}
			}

			if found {
				fmt.Printf("%s found, redirecting to : %s \r\n", uniqId, flink.Data.LongUrl)
				http.Redirect(w, r, flink.Data.LongUrl, http.StatusPermanentRedirect)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Printf("%s not found in shorted links\r\n", uniqId)
			fmt.Fprintln(w, "unique id not found")
		} else {
			if s.releasemode == false {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(s.shortlinks)
				return
			}
			fmt.Fprintln(w, "no unique id specifiq in URL : this is the root path")
		}
	}
}

func (s *Server) Run() {
	fmt.Printf("Server started on port :%d\r\n", s.port)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", rootHandler(s))
	muxWithMiddleware := NewMiddleware("browserUserAgent", mux)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), muxWithMiddleware)

}
