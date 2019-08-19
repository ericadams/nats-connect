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
	mux.HandleFunc("/version", versionHandler)
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func versionHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(version + "\n"))
}

func (s *server) connzHandler(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(s.db.Stats())
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(b)
}

// newServeMux builds a ServeMux and populates it with standard pprof handlers.
func (s *server) newServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	s.handle(mux)
	return mux
}

// newServer constructs a server at addr with the standard pprof handlers.
func newServer(opts *option, db *sql.DB) *server {
	return &server{
		opts: opts,
		db:   db,
	}
}

// Run a standard http server at addr.
func (s *server) Run() error {
	addr := fmt.Sprintf(":%d", s.opts.Port)
	log.Printf("starting server on %s", addr)

	return http.ListenAndServe(addr, s.newServeMux())
}
