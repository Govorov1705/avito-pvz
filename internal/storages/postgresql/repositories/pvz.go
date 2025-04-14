package repositories

import (
	"context"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"
	"github.com/Govorov1705/avito-pvz/pkg/logger"
	"github.com/Govorov1705/avito-pvz/pkg/transactions"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PvzRepository struct {
	Pool *pgxpool.Pool
}

func NewPvzRepository(p *pgxpool.Pool) *PvzRepository {
	return &PvzRepository{
		Pool: p,
	}
}

func (r *PvzRepository) Add(ctx context.Context, dto *dtos.CreatePvzRequest) (*models.Pvz, error) {
	pvz := models.Pvz{}

	stmt := `
		INSERT INTO pvzs(city)
		VALUES ($1)
		RETURNING id, registration_date, city;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var row pgx.Row
	if ok {
		row = tx.QueryRow(ctx, stmt, dto.City)
	} else {
		row = r.Pool.QueryRow(ctx, stmt, dto.City)
	}

	err := row.Scan(
		&pvz.ID,
		&pvz.RegistrationDate,
		&pvz.City,
	)
	if err != nil {
		logger.Logger.Error("Error scanning row", zap.Error(err))
		return nil, err
	}

	return &pvz, nil
}

func (r *PvzRepository) ListWithPagination(ctx context.Context, page, limit int) ([]*models.Pvz, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, registration_date, city
		FROM pvzs
		LIMIT $1
		OFFSET $2;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var rows pgx.Rows
	var err error
	if ok {
		rows, err = tx.Query(ctx, query, limit, offset)
	} else {
		rows, err = r.Pool.Query(ctx, query, limit, offset)
	}
	if err != nil {
		logger.Logger.Error("Error during query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var pvzs []*models.Pvz
	for rows.Next() {
		var pvz models.Pvz
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			logger.Logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		pvzs = append(pvzs, &pvz)
	}

	return pvzs, nil
}

func (r *PvzRepository) List(ctx context.Context) ([]*models.Pvz, error) {
	query := `
		SELECT id, registration_date, city
		FROM pvzs;
	`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		logger.Logger.Error("Error during query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var pvzs []*models.Pvz
	for rows.Next() {
		var pvz models.Pvz
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			logger.Logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		pvzs = append(pvzs, &pvz)
	}

	return pvzs, nil
}
