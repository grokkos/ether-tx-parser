package server

import (
	"github.com/grokkos/ether-tx-parser/internal/api/http/handler"
	"net/http"
)

type Server struct {
	handler *handler.ParserHandler
}

func NewServer(handler *handler.ParserHandler) *Server {
	return &Server{handler: handler}
}

func (s *Server) SetupRoutes() {
	http.HandleFunc("/block", s.handler.GetCurrentBlock)
	http.HandleFunc("/subscribe", s.handler.Subscribe)
	http.HandleFunc("/transactions", s.handler.GetTransactions)
}
