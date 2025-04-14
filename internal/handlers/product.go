package handlers

import (
	"errors"
	"net/http"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/services"
	"github.com/Govorov1705/avito-pvz/pkg/errs"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductService *services.ProductService
}

func NewProductHandler(s *services.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: s,
	}
}

func (h *ProductHandler) AddProductToOpenReception(c *gin.Context) {
	addProductRequest := dtos.AddProductRequest{}
	err := c.ShouldBindJSON(&addProductRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Error{Message: err.Error()})
		return
	}

	product, err := h.ProductService.Add(c.Request.Context(), &addProductRequest)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.JSON(http.StatusBadRequest, dtos.Error{Message: "Нет активной приемки"})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.Error{Message: "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusCreated, product)
}
