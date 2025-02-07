package parser

import (
	"encoding/json"
	"fmt"
	"github.com/grokkos/ether-tx-parser/internal/domain/entity"
	"github.com/grokkos/ether-tx-parser/internal/domain/repository"
	"github.com/grokkos/ether-tx-parser/pkg/errors"
	"github.com/grokkos/ether-tx-parser/pkg/logger"
	"go.uber.org/zap"
)

type Service struct {
	store  repository.Store
	client repository.EthereumClient
	logger *zap.Logger
}

func NewService(store repository.Store, client repository.EthereumClient) *Service {
	return &Service{
		store:  store,
		client: client,
		logger: logger.GetLogger(),
	}
}

func (s *Service) GetCurrentBlock() int {
	return s.store.GetCurrentBlock()
}

func (s *Service) Subscribe(address string) bool {
	// Validate Ethereum address format
	if len(address) != 42 || address[:2] != "0x" {
		s.logger.Warn("Invalid ethereum address format",
			zap.String("address", address),
		)
		return false
	}

	s.logger.Info("Subscribing to address",
		zap.String("address", address),
	)
	return s.store.Subscribe(address)
}

func (s *Service) GetTransactions(address string) []entity.Transaction {
	s.logger.Debug("Retrieving transactions",
		zap.String("address", address),
	)
	return s.store.GetTransactions(address)
}

type Block struct {
	Transactions []struct {
		Hash  string `json:"hash"`
		From  string `json:"from"`
		To    string `json:"to"`
		Value string `json:"value"`
	} `json:"transactions"`
}

func (s *Service) ParseBlocks() error {
	// Get latest block number
	response, err := s.client.MakeRPCCall("eth_blockNumber", []interface{}{})
	if err != nil {
		s.logger.Error("Failed to get latest block number",
			zap.Error(err),
		)
		return errors.NewEthereumError("failed to get latest block number", err)
	}

	blockNumberStr, ok := response.Result.(string)
	if !ok {
		s.logger.Error("Invalid block number format",
			zap.Any("response", response),
		)
		return errors.NewValidationError("invalid block number format", nil)
	}

	var latestBlock int
	fmt.Sscanf(blockNumberStr, "0x%x", &latestBlock)

	currentBlock := s.store.GetCurrentBlock()
	if currentBlock == 0 {
		currentBlock = latestBlock - 10
	}

	s.logger.Info("Starting block processing",
		zap.Int("current_block", currentBlock),
		zap.Int("latest_block", latestBlock),
	)

	// Process blocks
	for blockNum := currentBlock + 1; blockNum <= latestBlock; blockNum++ {
		if err := s.processBlock(blockNum); err != nil {
			s.logger.Error("Failed to process block",
				zap.Int("block_number", blockNum),
				zap.Error(err),
			)
			return errors.NewEthereumError(fmt.Sprintf("failed to process block %d", blockNum), err)
		}

		s.store.SetCurrentBlock(blockNum)
		s.logger.Debug("Processed block successfully",
			zap.Int("block_number", blockNum),
		)
	}

	return nil
}

func (s *Service) processBlock(blockNum int) error {
	blockResponse, err := s.client.MakeRPCCall("eth_getBlockByNumber",
		[]interface{}{fmt.Sprintf("0x%x", blockNum), true})
	if err != nil {
		return errors.NewEthereumError("failed to get block", err)
	}

	blockData, err := json.Marshal(blockResponse.Result)
	if err != nil {
		return errors.NewUnexpectedError("error marshaling block data", err)
	}

	var block Block
	if err := json.Unmarshal(blockData, &block); err != nil {
		return errors.NewValidationError("error unmarshaling block", err)
	}

	for _, tx := range block.Transactions {
		if s.store.IsSubscribed(tx.From) || s.store.IsSubscribed(tx.To) {
			s.logger.Debug("Found relevant transaction",
				zap.String("hash", tx.Hash),
				zap.String("from", tx.From),
				zap.String("to", tx.To),
			)

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
