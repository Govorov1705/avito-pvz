package middleware

import (
	"net/http"
	"strings"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"
	"github.com/Govorov1705/avito-pvz/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Auth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Error{Message: "Отсутствутет токен в заголовке Authorization"})
		return
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Error{Message: "Некорректный заголовок Authorization"})
		return
	}

	token := authHeaderParts[1]
	claims, err := jwt.ValidateJWT(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Error{Message: "Некорректный токен"})
		return
	}

	userIdStr := claims["sub"].(string)
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Error{Message: "Некорректный токен"})
		return
	}

	userRoleStr := claims["role"].(string)
	userRole := models.UserRole(userRoleStr)

	c.Set("userId", userId)
	c.Set("userRole", userRole)
	c.Next()
}
