package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
)

type server struct {
	opts *option
}

// handle adds standard http handlers to mux.
func handle(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/healthz", healthHandler)
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// newServeMux builds a ServeMux and populates it with standard pprof handlers.
func newServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	handle(mux)
	return mux
}

// newServer constructs a server at addr with the standard pprof handlers.
func newServer(opts *option) *server {
	return &server{
		opts: opts,
	}
}

// Run a standard http server at addr.
func (s *server) Run() error {
	addr := fmt.Sprintf(":%d", s.opts.Port)
	log.Printf("starting server on %s", addr)

	return http.ListenAndServe(addr, newServeMux())
}
