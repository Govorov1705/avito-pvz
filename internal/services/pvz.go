package services

import (
	"context"
	"strconv"
	"time"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/metrics"
	"github.com/Govorov1705/avito-pvz/internal/models"
	"github.com/Govorov1705/avito-pvz/internal/repositories"
	"github.com/Govorov1705/avito-pvz/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PvzService struct {
	PvzRepository       repositories.PvzRepository
	ReceptionRepository repositories.ReceptionRepository
	ProductRepository   repositories.ProductRepository
}

func NewPvzService(
	pvzr repositories.PvzRepository,
	rr repositories.ReceptionRepository,
	pr repositories.ProductRepository,
) *PvzService {
	return &PvzService{
		PvzRepository:       pvzr,
		ReceptionRepository: rr,
		ProductRepository:   pr,
	}
}

func (s *PvzService) CreatePvz(ctx context.Context, dto *dtos.CreatePvzRequest) (*models.Pvz, error) {
	pvz, err := s.PvzRepository.Add(ctx, dto)
	if err != nil {
		return nil, err
	}

	metrics.PvzCreatedTotal.Inc()

	return pvz, nil
}

func (s *PvzService) PvzInfo(ctx context.Context, dto *dtos.PvzInfoRequest) (dtos.PvzInfoReponse, error) {
	startDate, err := time.Parse(time.RFC3339, dto.StartDateStr)
	if err != nil {
		logger.Logger.Error("Error parsing startDate", zap.Error(err))
		return nil, err
	}

	endDate, err := time.Parse(time.RFC3339, dto.EndDateStr)
	if err != nil {
		logger.Logger.Error("Error parsing endDate", zap.Error(err))
		return nil, err
	}

	page, err := strconv.Atoi(dto.PageStr)
	if err != nil {
		logger.Logger.Error("Error parsing page", zap.Error(err))
		return nil, err
	}

	limit, err := strconv.Atoi(dto.LimitStr)
	if err != nil {
		logger.Logger.Error("Error parsing limit", zap.Error(err))
		return nil, err
	}

	pvzs, err := s.PvzRepository.ListWithPagination(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	pvzIds := make([]*uuid.UUID, len(pvzs))
	for i, pvz := range pvzs {
		pvzIds[i] = &pvz.ID
	}

	receptions, err := s.ReceptionRepository.ListByPvzIdsWithDateRange(
		ctx,
		pvzIds,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, err
	}

	pvzToReceptionMap := make(map[uuid.UUID][]*models.Reception)
	receptionIds := make([]*uuid.UUID, len(receptions))
	for i, reception := range receptions {
		pvzToReceptionMap[reception.PvzId] = append(pvzToReceptionMap[reception.PvzId], reception)
		receptionIds[i] = &reception.ID
	}

	products, err := s.ProductRepository.ListByReceptionIds(ctx, receptionIds)
	if err != nil {
		return nil, err
	}

	receptionToProductMap := make(map[uuid.UUID][]*models.Product)
	for _, product := range products {
		receptionToProductMap[product.ReceptionId] = append(receptionToProductMap[product.ReceptionId], product)
	}

	pvzInfoResponse := make(dtos.PvzInfoReponse, len(pvzs))
	for i, pvz := range pvzs {
		receptionsWithProducts := []dtos.ReceptionWithProducts{}

		pvzReceptions := pvzToReceptionMap[pvz.ID]
		for _, pvzReception := range pvzReceptions {
			receptionWithProducts := dtos.ReceptionWithProducts{
				Reception: *pvzReception,
				Products:  receptionToProductMap[pvzReception.ID],
			}
			receptionsWithProducts = append(receptionsWithProducts, receptionWithProducts)
		}

		pvzInfoResponse[i] = &dtos.PvzInfo{
			Pvz:        *pvz,
			Receptions: receptionsWithProducts,
		}

	}

	return pvzInfoResponse, nil
}
