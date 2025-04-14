package middleware

import (
	"net/http"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"

	"github.com/gin-gonic/gin"
)

func EmployeeOnly(c *gin.Context) {
	userRoleAny, _ := c.Get("userRole")
	userRole, _ := userRoleAny.(models.UserRole)

	if userRole != models.Employee {
		c.AbortWithStatusJSON(http.StatusForbidden, dtos.Error{Message: "Некорректная роль"})
		return
	}

	c.Next()
}
