package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func New() *Server {
	s := &Server{Router: mux.NewRouter()}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.Router.HandleFunc("/health", healthCheckHandler).Methods("GET")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
