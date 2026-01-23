package postgres

import (
	"context"
	"database/sql"
	"time"

	"accounting/internal/domain/entity"
	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	repoEntity "accounting/internal/repository/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}

// Mapper: Domain Entity -> Repository Entity
func toRepoUser(user *entity.User) *repoEntity.User {
	return &repoEntity.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// Mapper: Repository Entity -> Domain Entity
func toDomainUser(dbUser *repoEntity.User) *entity.User {
	return &entity.User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	dbUser := toRepoUser(user)

	// Set timestamps at repository layer
	now := time.Now()
	dbUser.CreatedAt = now
	dbUser.UpdatedAt = now

	query := `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		dbUser.ID,
		dbUser.Name,
		dbUser.Email,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Update domain entity with timestamps
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
SELECT id, name, email, created_at, updated_at
FROM users
WHERE id = $1
`

	var dbUser repoEntity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dbUser.ID,
		&dbUser.Name,
		&dbUser.Email,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toDomainUser(&dbUser), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
SELECT id, name, email, created_at, updated_at
FROM users
WHERE email = $1
`

	var dbUser repoEntity.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&dbUser.ID,
		&dbUser.Name,
		&dbUser.Email,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toDomainUser(&dbUser), nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	dbUser := toRepoUser(user)

	// Set updated timestamp at repository layer
	dbUser.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET name = $2, email = $3, updated_at = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		dbUser.ID,
		dbUser.Name,
		dbUser.Email,
		dbUser.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domainerrors.NewErrNotFound("user", user.ID)
	}

	// Update domain entity with new timestamp
	user.UpdatedAt = dbUser.UpdatedAt

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domainerrors.NewErrNotFound("user", id)
	}

	return nil
}

// Compile-time interface check
var _ interfaces.UserRepository = (*UserRepository)(nil)
