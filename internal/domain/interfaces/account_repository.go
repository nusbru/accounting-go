package interfaces

import (
	"accounting/internal/domain/entity"
	"context"
)

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	GetByID(ctx context.Context, id string) (*entity.Account, error)
	ListByUserID(ctx context.Context, userID string) ([]*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id string) error
}
