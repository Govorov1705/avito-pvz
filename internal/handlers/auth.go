package handlers

import (
	"errors"
	"net/http"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/services"
	"github.com/Govorov1705/avito-pvz/pkg/errs"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserService *services.UserService
}

func NewAuthHandler(s *services.UserService) *AuthHandler {
	return &AuthHandler{
		UserService: s,
	}
}

func (h *AuthHandler) DummyLogin(c *gin.Context) {
	dummyLoginRequest := dtos.DummyLoginRequest{}
	err := c.ShouldBindJSON(&dummyLoginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Error{Message: err.Error()})
		return
	}

	token, err := h.UserService.DummyLogin(c.Request.Context(), &dummyLoginRequest)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *AuthHandler) Register(c *gin.Context) {
	registerRequest := dtos.RegisterRequest{}
	err := c.ShouldBindJSON(&registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Error{Message: err.Error()})
		return
	}

	registerResponse, err := h.UserService.Register(c.Request.Context(), &registerRequest)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			c.JSON(http.StatusBadRequest, dtos.Error{Message: "Пользователь с таким email уже существует"})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.Error{Message: "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusCreated, registerResponse)
}

func (h *AuthHandler) Login(c *gin.Context) {
	loginRequest := dtos.LoginRequest{}
	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Error{Message: err.Error()})
		return
	}

	token, err := h.UserService.Login(c.Request.Context(), &loginRequest)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) || errors.Is(err, errs.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, dtos.Error{Message: "Неправильный email или пароль"})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.Error{Message: "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, token)
}
