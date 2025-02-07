package repository

import "github.com/grokkos/ether-tx-parser/internal/domain/entity"

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []entity.Transaction
}
