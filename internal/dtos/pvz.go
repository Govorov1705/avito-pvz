package dtos

import "github.com/Govorov1705/avito-pvz/internal/models"

type CreatePvzRequest struct {
	City string `json:"city" binding:"required,oneof=Москва Санкт-Петербург Казань"`
}

type PvzInfoRequest struct {
	StartDateStr string
	EndDateStr   string
	PageStr      string
	LimitStr     string
}

type ReceptionWithProducts struct {
	Reception models.Reception  `json:"reception"`
	Products  []*models.Product `json:"products"`
}

type PvzInfo struct {
	Pvz        models.Pvz              `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions"`
}

type PvzInfoReponse []*PvzInfo
