package parser

import (
	"encoding/json"
	"fmt"
	"github.com/grokkos/ether-tx-parser/internal/domain/entity"
	"github.com/grokkos/ether-tx-parser/internal/domain/repository"
)

// Service orchestrates the blockchain parsing operations.
type Service struct {
	store  repository.Store
	client repository.EthereumClient
}

func NewService(store repository.Store, client repository.EthereumClient) *Service {
	return &Service{
		store:  store,
		client: client,
	}
}

// GetCurrentBlock returns the last processed block number.
func (s *Service) GetCurrentBlock() int {
	return s.store.GetCurrentBlock()
}

// Subscribe adds an Ethereum address to the watch list, returns false if the address format is invalid.
func (s *Service) Subscribe(address string) bool {
	if len(address) != 42 || address[:2] != "0x" {
		return false
	}
	return s.store.Subscribe(address)
}

// GetTransactions returns all transactions for a given address.
func (s *Service) GetTransactions(address string) []entity.Transaction {
	return s.store.GetTransactions(address)
}

// Block represents the structure of an Ethereum block's transaction data.
type Block struct {
	Transactions []struct {
		Hash  string `json:"hash"`
		From  string `json:"from"`
		To    string `json:"to"`
		Value string `json:"value"`
	} `json:"transactions"`
}

// ParseBlocks processes new blocks from the Ethereum blockchain.
// It starts from the last processed block and continues to the latest block.
func (s *Service) ParseBlocks() error {
	// Get the latest block number from Ethereum
	response, err := s.client.MakeRPCCall("eth_blockNumber", []interface{}{})
	if err != nil {
		return fmt.Errorf("failed to get latest block: %v", err)
	}

	var latestBlock int
	fmt.Sscanf(response.Result, "0x%x", &latestBlock)

	currentBlock := s.store.GetCurrentBlock()
	if currentBlock == 0 {
		// If starting fresh, begin 10 blocks back
		currentBlock = latestBlock - 10
	}

	// Process each block from current to latest
	for blockNum := currentBlock + 1; blockNum <= latestBlock; blockNum++ {
		blockResponse, err := s.client.MakeRPCCall("eth_getBlockByNumber",
			[]interface{}{fmt.Sprintf("0x%x", blockNum), true})
		if err != nil {
			return fmt.Errorf("failed to get block %d: %v", blockNum, err)
		}

		if err := s.processBlock(blockNum, blockResponse.Result); err != nil {
			return fmt.Errorf("error processing block %d: %v", blockNum, err)
		}

		s.store.SetCurrentBlock(blockNum)
	}

	return nil
}

// processBlock handles the individual block data, extracting and storing relevant transactions.
func (s *Service) processBlock(blockNum int, blockData string) error {
	var block Block
	if err := json.Unmarshal([]byte(blockData), &block); err != nil {
		return fmt.Errorf("error unmarshaling block: %v", err)
	}

	for _, tx := range block.Transactions {
		// Only process transactions involving subscribed addresses
		if s.store.IsSubscribed(tx.From) || s.store.IsSubscribed(tx.To) {
			transaction := entity.Transaction{
				Hash:        tx.Hash,
				From:        tx.From,
				To:          tx.To,
				Value:       tx.Value,
				BlockNumber: blockNum,
			}
			s.store.AddTransaction(transaction)
		}
	}

	return nil
}
