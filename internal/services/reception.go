package services

import (
	"context"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/metrics"
	"github.com/Govorov1705/avito-pvz/internal/models"
	"github.com/Govorov1705/avito-pvz/internal/repositories"
	"github.com/Govorov1705/avito-pvz/internal/storages/postgresql"
	"github.com/Govorov1705/avito-pvz/pkg/logger"
	"github.com/Govorov1705/avito-pvz/pkg/transactions"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ReceptionService struct {
	ReceptionRepository repositories.ReceptionRepository
	DB                  postgresql.TxStarter
}

func NewReceptionService(
	r repositories.ReceptionRepository,
	DB postgresql.TxStarter,
) *ReceptionService {
	return &ReceptionService{
		ReceptionRepository: r,
		DB:                  DB,
	}
}

func (s *ReceptionService) CreateReception(ctx context.Context, dto *dtos.CreateReceptionRequest) (*models.Reception, error) {
	reception, err := s.ReceptionRepository.Add(ctx, dto)
	if err != nil {
		return nil, err
	}

	metrics.ReceptionsCreatedTotal.Inc()

	return reception, nil
}

func (s *ReceptionService) CloseReception(ctx context.Context, pvzId uuid.UUID) (reception *models.Reception, err error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		logger.Logger.Error("Error starting transaction", zap.Error(err))
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	ctx = transactions.PutTxIntoContext(ctx, tx)

	reception, err = s.ReceptionRepository.GetOpenByPvzId(ctx, pvzId)
	if err != nil {
		return nil, err
	}

	return s.ReceptionRepository.Close(ctx, reception.ID)
}
