package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hex4coder/go-url-shortener/models"
	"github.com/hex4coder/go-url-shortener/utils"
)

type ServerOpts struct {
	releasemode bool
	ssl         bool
	verbose     bool
	domain      string
	shortener   utils.ShortenerI
	shortlinks  []*models.ShortLink
	port        int
}

func WithShortener(shortener utils.ShortenerI) ServerFuncOpt {
	return func(f *ServerOpts) {
		f.shortener = shortener
	}
}
func WithDomain(domain string) ServerFuncOpt {
	return func(f *ServerOpts) {
		f.domain = domain
	}
}
func WithSSL(f *ServerOpts) {
	f.ssl = true
}

func WithReleaseMode(f *ServerOpts) {
	f.releasemode = true
}
func WithVerbose(f *ServerOpts) {
	f.verbose = true
}
func WithPort(port int) ServerFuncOpt {
	return func(f *ServerOpts) {
		f.port = port
	}
}

type Server struct {
	ServerOpts
}

type ServerFuncOpt func(*ServerOpts)

func DefaultServerOpts() ServerOpts {
	return ServerOpts{
		releasemode: false,
		verbose:     false,
		domain:      "he.i",
		ssl:         false,
		port:        1996,
	}
}
func NewServer(opts ...ServerFuncOpt) *Server {
	d := DefaultServerOpts()

	for _, opt := range opts {
		opt(&d)
	}
	return &Server{d}
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
			domain = fmt.Sprintf("192.168.77.205:%d", s.port)
		}

		link.ShortUrl = fmt.Sprintf("%s://%s/%s", basehttp, domain, link.UniqueCode)
		link.QrImageUrl = utils.EncodeURLToImageBase64(link.ShortUrl)
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
	//muxWithMiddleware := NewMiddleware("browserUserAgent", mux)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", s.port), mux))

}
