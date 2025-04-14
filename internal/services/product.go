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

type ProductService struct {
	ProductRepository   repositories.ProductRepository
	ReceptionRepository repositories.ReceptionRepository
	DB                  postgresql.TxStarter
}

func NewProductService(
	pr repositories.ProductRepository,
	rr repositories.ReceptionRepository,
	DB postgresql.TxStarter,
) *ProductService {
	return &ProductService{
		ProductRepository:   pr,
		ReceptionRepository: rr,
		DB:                  DB,
	}
}

func (s *ProductService) Add(ctx context.Context, dto *dtos.AddProductRequest) (product *models.Product, err error) {
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

	reception, err := s.ReceptionRepository.GetOpenByPvzId(ctx, dto.PvzId)
	if err != nil {
		return nil, err
	}

	product, err = s.ProductRepository.Add(ctx, dto.Type, reception.ID)
	if err != nil {
		return nil, err
	}

	metrics.ProductsAddedTotal.Inc()

	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, pvzId uuid.UUID) (err error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		logger.Logger.Error("Error starting transaction", zap.Error(err))
		return err
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

	reception, err := s.ReceptionRepository.GetOpenByPvzId(ctx, pvzId)
	if err != nil {
		return err
	}

	product, err := s.ProductRepository.GetLastByReceptionId(ctx, reception.ID)
	if err != nil {
		return err
	}

	return s.ProductRepository.Delete(ctx, product.ID)
}
