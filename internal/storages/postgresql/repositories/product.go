package repositories

import (
	"context"
	"errors"

	"github.com/Govorov1705/avito-pvz/internal/models"
	"github.com/Govorov1705/avito-pvz/pkg/errs"
	"github.com/Govorov1705/avito-pvz/pkg/logger"
	"github.com/Govorov1705/avito-pvz/pkg/transactions"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ProductRepository struct {
	Pool *pgxpool.Pool
}

func NewProductRepository(p *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		Pool: p,
	}
}

func (r *ProductRepository) Add(ctx context.Context, productType models.ProductType, receptionId uuid.UUID) (*models.Product, error) {
	product := models.Product{}

	stmt := `
		INSERT INTO products(type, reception_id)
		VALUES($1, $2)
		RETURNING id, date_time, type, reception_id;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var row pgx.Row
	if ok {
		row = tx.QueryRow(ctx, stmt, productType, receptionId)
	} else {
		row = r.Pool.QueryRow(ctx, stmt, productType, receptionId)
	}

	err := row.Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionId,
	)
	if err != nil {
		logger.Logger.Error("Error scanning row", zap.Error(err))
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetLastByReceptionId(ctx context.Context, receptionId uuid.UUID) (*models.Product, error) {
	product := models.Product{}

	query := `
		SELECT id, date_time, type, reception_id
		FROM products
		WHERE reception_id = $1
		ORDER BY date_time DESC
		LIMIT 1;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var row pgx.Row
	if ok {
		row = tx.QueryRow(ctx, query, receptionId)
	} else {
		row = r.Pool.QueryRow(ctx, query, receptionId)
	}

	err := row.Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		logger.Logger.Error("Error scanning row", zap.Error(err))
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) Delete(ctx context.Context, productId uuid.UUID) error {
	stmt := `
		DELETE FROM products
		WHERE id = $1;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var (
		tag pgconn.CommandTag
		err error
	)
	if ok {
		tag, err = tx.Exec(ctx, stmt, productId)
	} else {
		tag, err = r.Pool.Exec(ctx, stmt, productId)
	}
	if err != nil {
		logger.Logger.Error("Error executing statement", zap.Error(err))
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *ProductRepository) ListByReceptionIds(ctx context.Context, receptionIds []*uuid.UUID) ([]*models.Product, error) {
	query := `
		SELECT id, date_time, type, reception_id
		FROM products
		WHERE reception_id = ANY($1);
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var rows pgx.Rows
	var err error
	if ok {
		rows, err = tx.Query(ctx, query, receptionIds)
	} else {
		rows, err = r.Pool.Query(ctx, query, receptionIds)
	}
	if err != nil {
		logger.Logger.Error("Error during query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.DateTime,
			&product.Type,
			&product.ReceptionId,
		); err != nil {
			logger.Logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}
