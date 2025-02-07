package parser

import (
	"encoding/json"
	"fmt"
	"github.com/grokkos/ether-tx-parser/internal/domain/entity"
	"github.com/grokkos/ether-tx-parser/pkg/ethereum"
	"testing"
)

// MockStore is our test implementation of the Store interface
type MockStore struct {
	currentBlock int
	subscribers  map[string]bool
	transactions map[string][]entity.Transaction
}

func NewMockStore() *MockStore {
	return &MockStore{
		subscribers:  make(map[string]bool),
		transactions: make(map[string][]entity.Transaction),
	}
}

// Implement all Store interface methods
func (m *MockStore) GetCurrentBlock() int {
	return m.currentBlock
}

func (m *MockStore) SetCurrentBlock(block int) {
	m.currentBlock = block
}

func (m *MockStore) Subscribe(address string) bool {
	m.subscribers[address] = true
	return true
}

func (m *MockStore) IsSubscribed(address string) bool {
	return m.subscribers[address]
}

func (m *MockStore) GetTransactions(address string) []entity.Transaction {
	return m.transactions[address]
}

func (m *MockStore) AddTransaction(tx entity.Transaction) {
	if m.subscribers[tx.From] {
		m.transactions[tx.From] = append(m.transactions[tx.From], tx)
	}
	if m.subscribers[tx.To] {
		m.transactions[tx.To] = append(m.transactions[tx.To], tx)
	}
}

// MockEthereumClient is our test implementation of the EthereumClient interface
type MockEthereumClient struct {
	blockNumber    string
	blockResponses map[string]string
	shouldFail     bool
}

func (m *MockEthereumClient) MakeRPCCall(method string, params []interface{}) (*ethereum.JSONRPCResponse, error) {
	if m.shouldFail {
		return nil, fmt.Errorf("mock error")
	}

	if method == "eth_blockNumber" {
		return &ethereum.JSONRPCResponse{
			Result: m.blockNumber,
		}, nil
	}

	if method == "eth_getBlockByNumber" {
		blockNum := params[0].(string)
		if response, ok := m.blockResponses[blockNum]; ok {
			// Parse the JSON string into a map[string]interface{}
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(response), &result); err != nil {
				return nil, fmt.Errorf("error parsing mock response: %v", err)
			}
			return &ethereum.JSONRPCResponse{
				Result: result,
			}, nil
		}
	}

	return nil, fmt.Errorf("unexpected method: %s", method)
}

func TestService_Subscribe(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "valid ethereum address",
			address: "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			want:    true,
		},
		{
			name:    "invalid address - too short",
			address: "0x742d35",
			want:    false,
		},
		{
			name:    "invalid address - no prefix",
			address: "742d35Cc6634C0532925a3b844Bc454e4438f44e",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMockStore()
			client := &MockEthereumClient{}
			service := NewService(store, client)

			got := service.Subscribe(tt.address)
			if got != tt.want {
				t.Errorf("Service.Subscribe() = %v, want %v", got, tt.want)
			}

			// If subscription should succeed, verify address is stored
			if tt.want {
				if !store.IsSubscribed(tt.address) {
					t.Errorf("Address was not stored in subscribers")
				}
			}
		})
	}
}

func TestService_ParseBlocks(t *testing.T) {
	// Create a mock block response
	blockJSON := `{
        "transactions": [
            {
                "hash": "0x123",
                "from": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
                "to": "0x842d35Cc6634C0532925a3b844Bc454e4438f44f",
                "value": "0x2386f26fc10000"
            }
        ]
    }`

	tests := []struct {
		name         string
		currentBlock int
		latestBlock  string
		blockData    map[string]string
		subscribed   string
		shouldFail   bool
		wantErr      bool
		wantTxCount  int
	}{
		{
			name:         "successful parse",
			currentBlock: 0x1b3,
			latestBlock:  "0x1b4",
			blockData: map[string]string{
				"0x1b4": blockJSON,
			},
			subscribed:  "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			shouldFail:  false,
			wantErr:     false,
			wantTxCount: 1,
		},
		{
			name:         "no new blocks",
			currentBlock: 0x1b4,
			latestBlock:  "0x1b4",
			blockData:    map[string]string{},
			subscribed:   "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			shouldFail:   false,
			wantErr:      false,
			wantTxCount:  0,
		},
		{
			name:         "client error",
			currentBlock: 0x1b3,
			latestBlock:  "0x1b4",
			blockData:    map[string]string{},
			shouldFail:   true,
			wantErr:      true,
			wantTxCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMockStore()
			store.SetCurrentBlock(tt.currentBlock)

			if tt.subscribed != "" {
				store.Subscribe(tt.subscribed)
			}

			client := &MockEthereumClient{
				blockNumber:    tt.latestBlock,
				blockResponses: tt.blockData,
				shouldFail:     tt.shouldFail,
			}

			service := NewService(store, client)
			err := service.ParseBlocks()

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ParseBlocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.subscribed != "" {
				txs := store.GetTransactions(tt.subscribed)
				if len(txs) != tt.wantTxCount {
					t.Errorf("Got %d transactions, want %d", len(txs), tt.wantTxCount)
				}
			}
		})
	}
}

func TestService_GetTransactions(t *testing.T) {
	store := NewMockStore()
	client := &MockEthereumClient{}
	service := NewService(store, client)

	address := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	store.Subscribe(address)

	// Add some test transactions
	testTx := entity.Transaction{
		Hash:        "0x123",
		From:        address,
		To:          "0x456",
		Value:       "0x2386f26fc10000",
		BlockNumber: 123,
	}
	store.AddTransaction(testTx)

	// Test getting transactions
	txs := service.GetTransactions(address)
	if len(txs) != 1 {
		t.Errorf("GetTransactions() returned %d transactions, want 1", len(txs))
	}

	if txs[0].Hash != testTx.Hash {
		t.Errorf("Transaction hash = %s, want %s", txs[0].Hash, testTx.Hash)
	}
}

func TestService_GetCurrentBlock(t *testing.T) {
	store := NewMockStore()
	client := &MockEthereumClient{}
	service := NewService(store, client)

	expectedBlock := 12345
	store.SetCurrentBlock(expectedBlock)

	if got := service.GetCurrentBlock(); got != expectedBlock {
		t.Errorf("GetCurrentBlock() = %v, want %v", got, expectedBlock)
	}
}
