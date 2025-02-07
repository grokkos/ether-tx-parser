package server

import (
	"github.com/grokkos/ether-tx-parser/internal/api/http/handler"
	"net/http"
)

type Server struct {
	handler *handler.ParserHandler
	mux     *http.ServeMux // Add this
}

func NewServer(handler *handler.ParserHandler) *Server {
	return &Server{
		handler: handler,
		mux:     http.NewServeMux(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) SetupRoutes() {
	s.mux.HandleFunc("/block", s.handler.GetCurrentBlock)
	s.mux.HandleFunc("/subscribe", s.handler.Subscribe)
	s.mux.HandleFunc("/transactions", s.handler.GetTransactions)
}
