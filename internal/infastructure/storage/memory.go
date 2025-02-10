package storage

import (
	"github.com/grokkos/ether-tx-parser/internal/domain/entity"
	"go.uber.org/zap"
	"strings"
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
		mutex:        &sync.RWMutex{},
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
	if s == nil || address == "" {
		return false
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.subscribers == nil {
		s.subscribers = make(map[string]bool)
	}
	// Normalize the address as without this we didn't match correctly in the processing
	address = strings.ToLower(address)
	s.subscribers[address] = true
	return true
}

func (s *MemoryStore) IsSubscribed(address string) bool {
	if s == nil || address == "" {
		return false
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.subscribers == nil {
		return false
	}

	address = strings.ToLower(address)
	return s.subscribers[address]
}

func (s *MemoryStore) AddTransaction(tx entity.Transaction) {
	if s == nil {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.transactions == nil {
		s.transactions = make(map[string][]entity.Transaction)
	}

	from := strings.ToLower(tx.From)
	to := strings.ToLower(tx.To)

	if s.subscribers[from] {
		s.transactions[from] = append(s.transactions[from], tx)
	}
	if s.subscribers[to] {
		s.transactions[to] = append(s.transactions[to], tx)
	}
}

func (s *MemoryStore) GetTransactions(address string) []entity.Transaction {
	if s == nil || address == "" {
		return []entity.Transaction{}
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.transactions == nil {
		return []entity.Transaction{}
	}

	address = strings.ToLower(address)
	if transactions, exists := s.transactions[address]; exists {
		return transactions
	}
	return []entity.Transaction{}
}
