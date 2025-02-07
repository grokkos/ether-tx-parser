package handler

import (
	"encoding/json"
	"github.com/grokkos/ether-tx-parser/internal/application/parser"
	"net/http"
)

type ParserHandler struct {
	service *parser.Service
}

func NewParserHandler(service *parser.Service) *ParserHandler {
	return &ParserHandler{service: service}
}

type SubscribeRequest struct {
	Address string `json:"address"`
}

func (h *ParserHandler) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := h.service.GetCurrentBlock()
	err := json.NewEncoder(w).Encode(map[string]int{"current_block": block})
	if err != nil {
		return
	}
}

func (h *ParserHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req SubscribeRequest

	// Check if method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	success := h.service.Subscribe(req.Address)
	err := json.NewEncoder(w).Encode(map[string]bool{"success": success})
	if err != nil {
		return
	}
}

func (h *ParserHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	transactions := h.service.GetTransactions(address)
	err := json.NewEncoder(w).Encode(transactions)
	if err != nil {
		return
	}
}
