package storage

import (
	"github.com/grokkos/ether-tx-parser/internal/domain/entity"
	"go.uber.org/zap"
	"sync"
)

type MemoryStore struct {
	currentBlock int
	subscribers  map[string]bool
	transactions map[string][]entity.Transaction
	mutex        *sync.RWMutex
	logger       *zap.Logger
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		subscribers:  make(map[string]bool),
		transactions: make(map[string][]entity.Transaction),
		mutex:        &sync.RWMutex{}, // Make sure mutex is initialized
	}
}

func (s *MemoryStore) GetCurrentBlock() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.currentBlock
}

func (s *MemoryStore) SetCurrentBlock(block int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.currentBlock = block
}

func (s *MemoryStore) Subscribe(address string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.subscribers[address] = true
	return true
}

func (s *MemoryStore) IsSubscribed(address string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.subscribers[address]
}

func (s *MemoryStore) GetTransactions(address string) []entity.Transaction {
	if s == nil {
		return []entity.Transaction{}
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Return empty slice instead of nil if no transactions
	if transactions, exists := s.transactions[address]; exists {
		return transactions
	}
	return []entity.Transaction{}
}

func (s *MemoryStore) AddTransaction(tx entity.Transaction) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.subscribers[tx.From] {
		s.transactions[tx.From] = append(s.transactions[tx.From], tx)
	}
	if s.subscribers[tx.To] {
		s.transactions[tx.To] = append(s.transactions[tx.To], tx)
	}
}
