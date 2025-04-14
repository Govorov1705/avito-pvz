package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"
	"github.com/Govorov1705/avito-pvz/internal/repositories/mocks"
	"github.com/Govorov1705/avito-pvz/internal/services"
	"github.com/Govorov1705/avito-pvz/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Register(t *testing.T) {
	type testCase struct {
		name        string
		dto         *dtos.RegisterRequest
		setupMock   func(r *mocks.MockUserRepository)
		expectError bool
	}

	tcs := []testCase{
		{
			name: "OK",
			dto: &dtos.RegisterRequest{
				Email:    "testuser@mail.com",
				Password: "testpassword",
				Role:     "employee",
			},
			setupMock: func(r *mocks.MockUserRepository) {
				r.On("GetByEmail", mock.Anything, "testuser@mail.com").
					Return(nil, errs.ErrNotFound)
				r.On("Add", mock.Anything, mock.Anything).
					Return(&models.User{
						ID:       uuid.New(),
						Email:    "testuser@mail.com",
						Password: "$2a$12$rXN14oKCZ9kzrgvatH2lT.C/c9K2Gk9ybqK8FuWcP498FepfvYema",
						Role:     "employee",
					}, nil)
			},
			expectError: false,
		},

		{
			name: "User already exists",
			dto: &dtos.RegisterRequest{
				Email:    "testuser@mail.com",
				Password: "testpassword",
				Role:     "employee",
			},
			setupMock: func(r *mocks.MockUserRepository) {
				r.On("GetByEmail", mock.Anything, "testuser@mail.com").
					Return(&models.User{
						ID:       uuid.New(),
						Email:    "testuser@mail.com",
						Password: "$2a$12$rXN14oKCZ9kzrgvatH2lT.C/c9K2Gk9ybqK8FuWcP498FepfvYema",
						Role:     "employee",
					}, nil)
			},
			expectError: true,
		},

		{
			name: "UserRepository GetByEmail error",
			dto: &dtos.RegisterRequest{
				Email:    "testuser@mail.com",
				Password: "testpassword",
				Role:     "employee",
			},
			setupMock: func(r *mocks.MockUserRepository) {
				r.On("GetByEmail", mock.Anything, "testuser@mail.com").
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},

		{
			name: "UserRepository Add error",
			dto: &dtos.RegisterRequest{
				Email:    "testuser@mail.com",
				Password: "testpassword",
				Role:     "employee",
			},
			setupMock: func(r *mocks.MockUserRepository) {
				r.On("GetByEmail", mock.Anything, "testuser@mail.com").
					Return(nil, errs.ErrNotFound)
				r.On("Add", mock.Anything, mock.Anything).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := mocks.NewMockUserRepository(t)

			tc.setupMock(mockUserRepo)

			userService := services.NewUserService(mockUserRepo)

			registerResponse, err := userService.Register(context.Background(), tc.dto)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, registerResponse)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, registerResponse)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	type testCase struct {
		name        string
		dto         *dtos.LoginRequest
		setupMock   func(r *mocks.MockUserRepository)
		expectError bool
	}

	tcs := []testCase{
		{
			name: "OK",
			dto: &dtos.LoginRequest{
				Email:    "testuser@mail.com",
				Password: "testpassword",
			},
			setupMock: func(r *mocks.MockUserRepository) {
				r.On("GetByEmail", mock.Anything, "testuser@mail.com").
					Return(&models.User{
						ID:       uuid.New(),
						Email:    "testuser@mail.com",
						Password: "$2a$12$rXN14oKCZ9kzrgvatH2lT.C/c9K2Gk9ybqK8FuWcP498FepfvYema",
						Role:     "employee",
					}, nil)
			},
			expectError: false,
		},

		{
			name: "Wrong credentials",
			dto: &dtos.LoginRequest{
				Email:    "testuser@mail.com",
				Password: "wrongpassword",
			},
			setupMock: func(r *mocks.MockUserRepository) {
				r.On("GetByEmail", mock.Anything, "testuser@mail.com").
					Return(&models.User{
						ID:       uuid.New(),
						Email:    "testuser@mail.com",
						Password: "$2a$12$rXN14oKCZ9kzrgvatH2lT.C/c9K2Gk9ybqK8FuWcP498FepfvYema",
						Role:     "employee",
					}, nil)
			},
			expectError: true,
		},

		{
			name: "Non existant user",
			dto: &dtos.LoginRequest{
				Email:    "testuser@mail.com",
				Password: "testpassword",
			},
			setupMock: func(r *mocks.MockUserRepository) {
				r.On("GetByEmail", mock.Anything, "testuser@mail.com").
					Return(nil, errs.ErrNotFound)
			},
			expectError: true,
		},

		{
			name: "UserRepo GetByEmail error",
			dto: &dtos.LoginRequest{
				Email:    "testuser@mail.com",
				Password: "testpassword",
			},
			setupMock: func(r *mocks.MockUserRepository) {
				r.On("GetByEmail", mock.Anything, "testuser@mail.com").
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := mocks.NewMockUserRepository(t)

			tc.setupMock(mockUserRepo)

			userService := services.NewUserService(mockUserRepo)

			token, err := userService.Login(context.Background(), tc.dto)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}
