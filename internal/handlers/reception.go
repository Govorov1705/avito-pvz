package handlers

import (
	"errors"
	"net/http"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/services"
	"github.com/Govorov1705/avito-pvz/pkg/errs"

	"github.com/gin-gonic/gin"
)

type ReceptionHandler struct {
	ReceptionService *services.ReceptionService
}

func NewReceptionHandler(s *services.ReceptionService) *ReceptionHandler {
	return &ReceptionHandler{
		ReceptionService: s,
	}
}

func (h *ReceptionHandler) CreateReception(c *gin.Context) {
	createReceptionRequest := dtos.CreateReceptionRequest{}
	err := c.ShouldBindJSON(&createReceptionRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Error{Message: err.Error()})
		return
	}

	reception, err := h.ReceptionService.CreateReception(c.Request.Context(), &createReceptionRequest)
	if err != nil {
		if errors.Is(err, errs.ErrUniqueConstraintViolation) {
			c.JSON(http.StatusBadRequest, dtos.Error{Message: "Уже существует открытая приемка"})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.Error{Message: "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusCreated, reception)
}
