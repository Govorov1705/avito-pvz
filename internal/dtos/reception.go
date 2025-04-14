package dtos

import "github.com/google/uuid"

type CreateReceptionRequest struct {
	PvzId uuid.UUID `json:"pvz_id" binding:"required,uuid4"`
}
