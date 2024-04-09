package web

import (
	"net/http"
)

type Server struct {
	http.Server
}

func NewServer(host string) Server {
	return Server{
		http.Server{
			Addr: host,
		},
	}
}

func (s *Server) Run(mux *http.ServeMux) error {
	return http.ListenAndServe(s.Addr, mux)
}
