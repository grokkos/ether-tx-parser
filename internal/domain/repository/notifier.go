package repository

import "github.com/grokkos/ether-tx-parser/internal/domain/entity"

// NotificationService interface defined but not implemented
// This serves as a hook point for future notification implementations
type NotificationService interface {
	NotifyTransaction(tx entity.Transaction) error
}
