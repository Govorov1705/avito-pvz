package repositories

import (
	"context"

	"github.com/Govorov1705/avito-pvz/internal/models"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Add(ctx context.Context, productType models.ProductType, receptionId uuid.UUID) (*models.Product, error)
	GetLastByReceptionId(ctx context.Context, receptionId uuid.UUID) (*models.Product, error)
	Delete(ctx context.Context, productId uuid.UUID) error
	ListByReceptionIds(ctx context.Context, receptionIds []*uuid.UUID) ([]*models.Product, error)
}
