package jwt

import (
	"errors"
	"fmt"

	"github.com/Govorov1705/avito-pvz/configs"
	"github.com/Govorov1705/avito-pvz/internal/models"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateJWT(userId uuid.UUID, role models.UserRole) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userId,
		"role": role,
		"iat":  now.Unix(),
		"exp":  now.Add(10 * 24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(configs.Cfg.SecretKey))
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.Cfg.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("invalid token claims type")
	}

	return claims, nil
}
