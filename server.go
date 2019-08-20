package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"

	_ "github.com/denisenkom/go-mssqldb"
)

type server struct {
	opts *option
	conn *Connector
}

// handle adds standard http handlers to mux.
func (s *server) newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.HandleFunc("/healthz", healthHandler)
	r.HandleFunc("/version", versionHandler)
	r.HandleFunc("/connectors", s.connectorsHandler)
	r.HandleFunc("/sources", s.sourcesHandler)

	return r
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func versionHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(version + "\n"))
}

func (s *server) connectorsHandler(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(s.conn.Stats())
	if err != nil {
		log.Printf("error:%v", err)
	}
	w.Write(b)
}

func (s *server) sourcesHandler(w http.ResponseWriter, req *http.Request) {
	if len(s.conn.Sources) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(s.conn.Sources)
	if err != nil {
		log.Printf("error:%v", err)
	}
	w.Write(b)
}

// newServer constructs a server at addr with the standard pprof handlers.
func newServer(opts *option, conn *Connector) *server {
	return &server{
		opts: opts,
		conn: conn,
	}
}

// Run a standard http server at addr.
func (s *server) Run() error {
	addr := fmt.Sprintf(":%d", s.opts.Port)
	log.Printf("starting server on %s", addr)

	return http.ListenAndServe(addr, s.newRouter())
}
