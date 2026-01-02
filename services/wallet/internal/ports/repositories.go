package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/wallet/internal/domain"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error)
	Update(ctx context.Context, wallet *domain.Wallet) error
	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error)
	GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*domain.Transaction, error)
	Update(ctx context.Context, tx *domain.Transaction) error
	CountByWalletID(ctx context.Context, walletID uuid.UUID) (int, error)
}

type PaymentMethodRepository interface {
	Create(ctx context.Context, pm *domain.PaymentMethod) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PaymentMethod, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.PaymentMethod, error)
	GetDefaultByUserID(ctx context.Context, userID uuid.UUID) (*domain.PaymentMethod, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SetDefault(ctx context.Context, userID, methodID uuid.UUID) error
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(tx Transaction) error) error
}

type Transaction interface {
	Wallets() WalletRepository
	Transactions() TransactionRepository
}
