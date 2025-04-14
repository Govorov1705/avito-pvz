package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
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

type ReceptionRepository struct {
	Pool *pgxpool.Pool
}

func NewReceptionRepository(p *pgxpool.Pool) *ReceptionRepository {
	return &ReceptionRepository{
		Pool: p,
	}
}

func (r *ReceptionRepository) Add(ctx context.Context, dto *dtos.CreateReceptionRequest) (*models.Reception, error) {
	reception := models.Reception{}

	stmt := `
		INSERT INTO receptions(pvz_id)
		VALUES($1)
		RETURNING id, date_time, pvz_id, status;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var row pgx.Row
	if ok {
		row = tx.QueryRow(ctx, stmt, dto.PvzId)
	} else {
		row = r.Pool.QueryRow(ctx, stmt, dto.PvzId)
	}

	err := row.Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PvzId,
		&reception.Status,
	)

	var pgErr *pgconn.PgError
	if err != nil {
		// Only one open reception allowed per pvz
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errs.ErrUniqueConstraintViolation
		}
		logger.Logger.Error("Error scanning row", zap.Error(err))
		return nil, err
	}

	return &reception, nil
}

func (r *ReceptionRepository) GetOpenByPvzId(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	reception := models.Reception{}

	query := `
		SELECT id, date_time, pvz_id, status
		FROM receptions
		WHERE pvz_id = $1 AND status = $2
		FOR UPDATE;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var row pgx.Row
	if ok {
		row = tx.QueryRow(ctx, query, pvzId, models.InProgress)
	} else {
		row = r.Pool.QueryRow(ctx, query, pvzId, models.InProgress)
	}

	err := row.Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PvzId,
		&reception.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		logger.Logger.Error("Error scanning row", zap.Error(err))
		return nil, err
	}

	return &reception, nil
}

func (r *ReceptionRepository) Close(ctx context.Context, receptionId uuid.UUID) (*models.Reception, error) {
	reception := models.Reception{}

	stmt := `
		UPDATE receptions
		SET status = $1
		WHERE id = $2 and status = $3
		RETURNING id, date_time, pvz_id, status;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var row pgx.Row
	if ok {
		row = tx.QueryRow(ctx, stmt, models.Close, receptionId, models.InProgress)
	} else {
		row = r.Pool.QueryRow(ctx, stmt, models.Close, receptionId, models.InProgress)
	}

	err := row.Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PvzId,
		&reception.Status,
	)
	if err != nil {
		logger.Logger.Error("Error scanning row", zap.Error(err))
		return nil, err
	}

	return &reception, nil
}

func (r *ReceptionRepository) ListByPvzIdsWithDateRange(ctx context.Context, pvzIds []*uuid.UUID, startDate, endDate time.Time) ([]*models.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status
		FROM receptions
		WHERE pvz_id = ANY($1) 
		AND date_time BETWEEN $2 AND $3;
	`

	tx, ok := transactions.GetTxFromContext(ctx)
	var rows pgx.Rows
	var err error
	if ok {
		rows, err = tx.Query(ctx, query, pvzIds, startDate, endDate)
	} else {
		rows, err = r.Pool.Query(ctx, query, pvzIds, startDate, endDate)
	}
	if err != nil {
		logger.Logger.Error("Error during query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var receptions []*models.Reception
	for rows.Next() {
		var reception models.Reception
		if err := rows.Scan(
			&reception.ID,
			&reception.DateTime,
			&reception.PvzId,
			&reception.Status,
		); err != nil {
			logger.Logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		receptions = append(receptions, &reception)
	}

	return receptions, nil
}
