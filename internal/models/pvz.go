package models

import (
	"time"

	"github.com/google/uuid"
)

type PvzCity string

const (
	Moscow          PvzCity = "Москва"
	SaintPetersburg PvzCity = "Санкт-Петербург"
	Kazan           PvzCity = "Казань"
)

type Pvz struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             PvzCity   `json:"city"`
}
