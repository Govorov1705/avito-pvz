package repositories

import (
	"context"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Add(ctx context.Context, dto *dtos.RegisterRequest) (*models.User, error)
	GetById(ctx context.Context, userId uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
