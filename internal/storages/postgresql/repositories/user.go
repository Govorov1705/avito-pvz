package repositories

import (
	"context"
	"errors"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"
	"github.com/Govorov1705/avito-pvz/pkg/errs"
	"github.com/Govorov1705/avito-pvz/pkg/logger"
	"github.com/Govorov1705/avito-pvz/pkg/transactions"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(p *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		Pool: p,
	}
}

func (r *UserRepository) Add(ctx context.Context, dto *dtos.RegisterRequest) (*models.User, error) {
	user := models.User{}

	stmt := `
		INSERT INTO users(
			email,
			password,
			role
		) values($1, $2, $3)
		RETURNING id, email, password, role;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var row pgx.Row
	if ok {
		row = tx.QueryRow(ctx, stmt, dto.Email, dto.Password, dto.Role)
	} else {
		row = r.Pool.QueryRow(ctx, stmt, dto.Email, dto.Password, dto.Role)
	}

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		logger.Logger.Error("Error scanning row", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetById(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	user := models.User{}

	query := `
		SELECT id, email, password, role
		FROM users
		WHERE id = $1;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var row pgx.Row
	if ok {
		row = tx.QueryRow(ctx, query, userId)
	} else {
		row = r.Pool.QueryRow(ctx, query, userId)
	}

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		logger.Logger.Error("Error scanning row", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := models.User{}

	query := `
		SELECT id, email, password, role
		FROM users
		WHERE email = $1;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var row pgx.Row
	if ok {
		row = tx.QueryRow(ctx, query, email)
	} else {
		row = r.Pool.QueryRow(ctx, query, email)
	}

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		logger.Logger.Error("Error scanning row", zap.Error(err))
		return nil, err
	}

	return &user, nil
}
