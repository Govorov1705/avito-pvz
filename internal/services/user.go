package services

import (
	"context"
	"errors"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"
	"github.com/Govorov1705/avito-pvz/internal/repositories"
	"github.com/Govorov1705/avito-pvz/pkg/errs"
	"github.com/Govorov1705/avito-pvz/pkg/jwt"
	"github.com/Govorov1705/avito-pvz/pkg/password"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) *UserService {
	return &UserService{
		UserRepository: r,
	}
}

func (s *UserService) DummyLogin(ctx context.Context, dto *dtos.DummyLoginRequest) (token string, err error) {
	return jwt.CreateJWT(uuid.New(), models.UserRole(dto.Role))
}

func (s *UserService) Register(ctx context.Context, dto *dtos.RegisterRequest) (*dtos.RegisterResponse, error) {
	_, err := s.UserRepository.GetByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			hashedPassword, err := password.HashPassword(dto.Password)
			if err != nil {
				return nil, err
			}
			dto.Password = hashedPassword

			user, err := s.UserRepository.Add(ctx, dto)
			if err != nil {
				return nil, err
			}
			registerResponse := dtos.RegisterResponse{
				ID:    user.ID,
				Email: user.Email,
				Role:  string(user.Role),
			}
			return &registerResponse, nil
		}

		return nil, err
	}

	return nil, errs.ErrAlreadyExists
}

func (s *UserService) Login(ctx context.Context, dto *dtos.LoginRequest) (token string, err error) {
	user, err := s.UserRepository.GetByEmail(ctx, dto.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return "", errs.ErrInvalidCredentials
	}

	token, err = jwt.CreateJWT(user.ID, user.Role)
	return token, err
}
