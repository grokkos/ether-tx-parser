package main

import (
	"github.com/grokkos/ether-tx-parser/internal/api/http/handler"
	"github.com/grokkos/ether-tx-parser/internal/api/http/server"
	"github.com/grokkos/ether-tx-parser/internal/application/parser"
	"github.com/grokkos/ether-tx-parser/internal/infastructure/ethereum"
	"github.com/grokkos/ether-tx-parser/internal/infastructure/storage"
	"log"
	"net/http"
	"time"
)

func main() {
	// Initialize dependencies
	store := storage.NewMemoryStore()
	client := ethereum.NewClient("https://ethereum-rpc.publicnode.com")
	service := parser.NewService(store, client)

	// Setup HTTP server
	parserHandler := handler.NewParserHandler(service)
	srv := server.NewServer(parserHandler)
	srv.SetupRoutes()

	// Start parsing blocks in a goroutine
	go func() {
		for {
			if err := service.ParseBlocks(); err != nil {
				log.Printf("Error parsing blocks: %v\n", err)
			}
			time.Sleep(15 * time.Second) // Poll for new blocks every 15 seconds
		}
	}()

	// Start HTTP server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
