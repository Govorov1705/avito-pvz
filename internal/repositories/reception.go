package repositories

import (
	"context"
	"time"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"

	"github.com/google/uuid"
)

type ReceptionRepository interface {
	Add(ctx context.Context, dto *dtos.CreateReceptionRequest) (*models.Reception, error)
	GetOpenByPvzId(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error)
	Close(ctx context.Context, receptionId uuid.UUID) (*models.Reception, error)
	ListByPvzIdsWithDateRange(ctx context.Context, pvzIds []*uuid.UUID, startDate, endDate time.Time) ([]*models.Reception, error)
}
