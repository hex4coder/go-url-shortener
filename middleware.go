package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Middleware struct {
	UserAgent string
	handler   http.Handler
}

// ServeHTTP implements http.Handler.
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("checking browser agent ....")
	userAgent := r.UserAgent()
	if strings.Contains(userAgent, m.UserAgent) {
		m.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "forbidden")
	}

}

func NewMiddleware(userAgentContains string, handle http.Handler) *Middleware {
	return &Middleware{
		UserAgent: userAgentContains,
		handler:   handle,
	}
}
