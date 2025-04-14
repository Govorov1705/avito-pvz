package dtos

import (
	"github.com/Govorov1705/avito-pvz/internal/models"

	"github.com/google/uuid"
)

type AddProductRequest struct {
	Type  models.ProductType `json:"type" binding:"required,oneof=электроника одежда обувь"`
	PvzId uuid.UUID          `json:"pvzId" binding:"required,uuid4"`
}
