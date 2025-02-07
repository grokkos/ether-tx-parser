package storage

import (
	"github.com/grokkos/ether-tx-parser/internal/domain/entity"
	"sync"
)

type MemoryStore struct {
	currentBlock int
	subscribers  map[string]bool
	transactions map[string][]entity.Transaction
	mutex        sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		subscribers:  make(map[string]bool),
		transactions: make(map[string][]entity.Transaction),
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
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.transactions[address]
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
