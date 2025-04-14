package handlers

import (
	"errors"
	"net/http"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/services"
	"github.com/Govorov1705/avito-pvz/pkg/errs"
	"github.com/Govorov1705/avito-pvz/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PvzHandler struct {
	PvzService       *services.PvzService
	ReceptionService *services.ReceptionService
	ProductService   *services.ProductService
}

func NewPvzHandler(
	pvzs *services.PvzService,
	rs *services.ReceptionService,
	ps *services.ProductService,

) *PvzHandler {
	return &PvzHandler{
		PvzService:       pvzs,
		ReceptionService: rs,
		ProductService:   ps,
	}
}

func (h *PvzHandler) CreatePvz(c *gin.Context) {
	createPvzRequest := dtos.CreatePvzRequest{}
	err := c.ShouldBindJSON(&createPvzRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Error{Message: err.Error()})
		return
	}

	pvz, err := h.PvzService.CreatePvz(c.Request.Context(), &createPvzRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Error{Message: "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusCreated, pvz)
}

func (h *PvzHandler) CloseLastReception(c *gin.Context) {
	pvzIdStr := c.Param("pvzId")
	pvdId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		logger.Logger.Error("Error parsing uuid from path param", zap.Error(err))
		c.JSON(http.StatusBadRequest, dtos.Error{Message: "Внутренняя ошибка сервера"})
		return
	}

	reception, err := h.ReceptionService.CloseReception(c.Request.Context(), pvdId)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.JSON(http.StatusBadRequest, dtos.Error{Message: "У данного ПВЗ нет открытой приемки"})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.Error{Message: "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, reception)
}

func (h *PvzHandler) DeleteLastReceptionProduct(c *gin.Context) {
	pvzIdStr := c.Param("pvzId")
	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		logger.Logger.Error("Error parsing uuid from path param", zap.Error(err))
		c.JSON(http.StatusBadRequest, dtos.Error{Message: "Внутренняя ошибка сервера"})
		return
	}

	err = h.ProductService.Delete(c.Request.Context(), pvzId)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.JSON(http.StatusBadRequest, dtos.Error{Message: "Неверный запрос, нет активной приемки или нет товаров для удаления"})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.Error{Message: "Внутренняя ошибка сервера"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *PvzHandler) PvzInfo(c *gin.Context) {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if startDateStr == "" || endDateStr == "" || pageStr == "" || limitStr == "" {
		c.JSON(
			http.StatusBadRequest,
			dtos.Error{Message: "В запросе должны присутстовать query-параметры startDate, endDate, page, limit"},
		)
		return
	}

	pvzInfoRequest := dtos.PvzInfoRequest{
		StartDateStr: startDateStr,
		EndDateStr:   endDateStr,
		PageStr:      pageStr,
		LimitStr:     limitStr,
	}

	pvzs, err := h.PvzService.PvzInfo(c.Request.Context(), &pvzInfoRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Error{Message: "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, pvzs)
}
