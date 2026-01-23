package interfaces

import (
	"accounting/internal/domain/entity"
	"context"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	ListByAccountID(ctx context.Context, accountID string) ([]*entity.Transaction, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	Delete(ctx context.Context, id string) error
}
