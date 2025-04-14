package models

import (
	"time"

	"github.com/google/uuid"
)

type ReceptionStatus string

const (
	InProgress ReceptionStatus = "in_progress"
	Close      ReceptionStatus = "close"
)

type Reception struct {
	ID       uuid.UUID       `json:"id"`
	DateTime time.Time       `json:"dateTime"`
	PvzId    uuid.UUID       `json:"pvzId"`
	Status   ReceptionStatus `json:"status"`
}
