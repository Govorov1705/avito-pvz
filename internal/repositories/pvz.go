package repositories

import (
	"context"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"
)

type PvzRepository interface {
	Add(ctx context.Context, dto *dtos.CreatePvzRequest) (*models.Pvz, error)
	ListWithPagination(ctx context.Context, page, limit int) ([]*models.Pvz, error)
	List(ctx context.Context) ([]*models.Pvz, error)
}
