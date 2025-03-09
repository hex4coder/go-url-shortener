package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

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

	fmt.Printf("[INFO] - Ditemukan %d data, lanjut memproses shortlinks....\r\n", len(datalinks))
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

func GetByUniqID(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Create browser parser only for Exam Browser
		fmt.Printf(
			"Incoming request from : r.UserAgent(): %v\n",
			r.UserAgent(),
		)

		// Get the full URL path
		uniqId := r.PathValue("id")
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
			fmt.Printf("%s found, redirecting to : %s \r\n", uniqId, flink.DataLink.LongUrl)
			http.Redirect(w, r, flink.DataLink.LongUrl, http.StatusPermanentRedirect)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("%s not found in shorted links\r\n", uniqId)
		fmt.Fprintln(w, "unique id not found", uniqId)
	}
}

func GetShortLinks(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(s.shortlinks) == 0 {
			fmt.Fprintln(w, "no shortlinks")
			return
		}
		// render html golang
		filepath := path.Join("views", "index.html")
		tmpl, err := template.ParseFiles(filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, s.shortlinks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) Run() {
	fmt.Printf("Server started on port :%d\r\n", s.port)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", GetShortLinks(s))
	mux.HandleFunc("GET /{id}", GetByUniqID(s))
	// muxWithMiddleware := NewMiddleware("browserUserAgent", mux)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", s.port), mux))
}
