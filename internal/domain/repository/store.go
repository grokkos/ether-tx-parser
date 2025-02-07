package repository

import (
	"github.com/grokkos/ether-tx-parser/internal/domain/entity"
	"github.com/grokkos/ether-tx-parser/pkg/ethereum"
)

type Store interface {
	GetCurrentBlock() int
	SetCurrentBlock(block int)
	Subscribe(address string) bool
	IsSubscribed(address string) bool
	GetTransactions(address string) []entity.Transaction
	AddTransaction(tx entity.Transaction)
}

// EthereumClient defines the interface for interacting with Ethereum nodes.
type EthereumClient interface {
	MakeRPCCall(method string, params []interface{}) (*ethereum.JSONRPCResponse, error)
}
