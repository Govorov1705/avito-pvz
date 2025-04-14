package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"
	"github.com/Govorov1705/avito-pvz/internal/repositories/mocks"
	"github.com/Govorov1705/avito-pvz/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPvzService_CreatePvz(t *testing.T) {
	type testCase struct {
		name        string
		dto         *dtos.CreatePvzRequest
		setupMock   func(r *mocks.MockPvzRepository)
		expectError bool
	}

	tcs := []testCase{
		{
			name: "OK",
			dto: &dtos.CreatePvzRequest{
				City: "Москва",
			},
			setupMock: func(r *mocks.MockPvzRepository) {
				r.On("Add", mock.Anything, &dtos.CreatePvzRequest{
					City: "Москва",
				}).Return(&models.Pvz{
					ID:               uuid.New(),
					RegistrationDate: time.Now(),
					City:             "Москва",
				}, nil)
			},
			expectError: false,
		},

		{
			name: "PvzRepository Add error",
			dto: &dtos.CreatePvzRequest{
				City: "Москва",
			},
			setupMock: func(r *mocks.MockPvzRepository) {
				r.On("Add", mock.Anything, &dtos.CreatePvzRequest{
					City: "Москва",
				}).Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockPvzRepo := mocks.NewMockPvzRepository(t)
			mockReceptionRepo := mocks.NewMockReceptionRepository(t)
			mockProductRepo := mocks.NewMockProductRepository(t)

			tc.setupMock(mockPvzRepo)

			pvzService := services.NewPvzService(mockPvzRepo, mockReceptionRepo, mockProductRepo)

			pvz, err := pvzService.CreatePvz(context.Background(), tc.dto)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, pvz)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, pvz)
			}

			mockPvzRepo.AssertExpectations(t)
		})
	}
}

func TestPvzService_PvzInfo(t *testing.T) {
	type testCase struct {
		name       string
		dto        *dtos.PvzInfoRequest
		setupMocks func(
			pvzr *mocks.MockPvzRepository,
			rr *mocks.MockReceptionRepository,
			pr *mocks.MockProductRepository,
		)
		expectError bool
	}

	startDateStr := "2025-01-01T15:00:00Z"
	endDateStr := "2025-05-11T10:00:00Z"

	tcs := []testCase{
		{
			name: "OK",
			dto: &dtos.PvzInfoRequest{
				StartDateStr: startDateStr,
				EndDateStr:   endDateStr,
				PageStr:      "1",
				LimitStr:     "1",
			},
			setupMocks: func(
				pvzr *mocks.MockPvzRepository,
				rr *mocks.MockReceptionRepository,
				pr *mocks.MockProductRepository,
			) {
				startDate, err := time.Parse(time.RFC3339, startDateStr)
				require.NoError(t, err)

				endDate, err := time.Parse(time.RFC3339, endDateStr)
				require.NoError(t, err)

				pvz := models.Pvz{
					ID:               uuid.New(),
					RegistrationDate: time.Date(2025, time.March, 18, 10, 0, 0, 0, time.UTC),
					City:             models.Moscow,
				}

				pvzr.On("ListWithPagination", mock.Anything, 1, 1).
					Return([]*models.Pvz{&pvz}, nil)

				reception := models.Reception{
					ID:       uuid.New(),
					DateTime: startDate,
					PvzId:    pvz.ID,
					Status:   models.Close,
				}

				rr.On("ListByPvzIdsWithDateRange", mock.Anything, []*uuid.UUID{&pvz.ID}, startDate, endDate).
					Return([]*models.Reception{&reception}, nil)

				product := models.Product{
					ID:          uuid.New(),
					DateTime:    startDate.Add(5 * time.Minute),
					Type:        models.Clothes,
					ReceptionId: reception.ID,
				}

				pr.On("ListByReceptionIds", mock.Anything, []*uuid.UUID{&reception.ID}).
					Return([]*models.Product{&product}, nil)
			},
			expectError: false,
		},

		{
			name: "PvzRepository ListWithPagination error",
			dto: &dtos.PvzInfoRequest{
				StartDateStr: startDateStr,
				EndDateStr:   endDateStr,
				PageStr:      "1",
				LimitStr:     "1",
			},
			setupMocks: func(
				pvzr *mocks.MockPvzRepository,
				rr *mocks.MockReceptionRepository,
				pr *mocks.MockProductRepository,
			) {
				pvzr.On("ListWithPagination", mock.Anything, 1, 1).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},

		{
			name: "ReceptionRepository ListByPvzIdsWithDateRange error",
			dto: &dtos.PvzInfoRequest{
				StartDateStr: startDateStr,
				EndDateStr:   endDateStr,
				PageStr:      "1",
				LimitStr:     "1",
			},
			setupMocks: func(
				pvzr *mocks.MockPvzRepository,
				rr *mocks.MockReceptionRepository,
				pr *mocks.MockProductRepository,
			) {
				startDate, err := time.Parse(time.RFC3339, startDateStr)
				require.NoError(t, err)

				endDate, err := time.Parse(time.RFC3339, endDateStr)
				require.NoError(t, err)

				pvz := models.Pvz{
					ID:               uuid.New(),
					RegistrationDate: time.Date(2025, time.March, 18, 10, 0, 0, 0, time.UTC),
					City:             models.Moscow,
				}

				pvzr.On("ListWithPagination", mock.Anything, 1, 1).
					Return([]*models.Pvz{&pvz}, nil)

				rr.On("ListByPvzIdsWithDateRange", mock.Anything, []*uuid.UUID{&pvz.ID}, startDate, endDate).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},

		{
			name: "ProductRepository ListByReceptionIds error",
			dto: &dtos.PvzInfoRequest{
				StartDateStr: startDateStr,
				EndDateStr:   endDateStr,
				PageStr:      "1",
				LimitStr:     "1",
			},
			setupMocks: func(
				pvzr *mocks.MockPvzRepository,
				rr *mocks.MockReceptionRepository,
				pr *mocks.MockProductRepository,
			) {
				startDate, err := time.Parse(time.RFC3339, startDateStr)
				require.NoError(t, err)

				endDate, err := time.Parse(time.RFC3339, endDateStr)
				require.NoError(t, err)

				pvz := models.Pvz{
					ID:               uuid.New(),
					RegistrationDate: time.Date(2025, time.March, 18, 10, 0, 0, 0, time.UTC),
					City:             models.Moscow,
				}

				pvzr.On("ListWithPagination", mock.Anything, 1, 1).
					Return([]*models.Pvz{&pvz}, nil)

				reception := models.Reception{
					ID:       uuid.New(),
					DateTime: startDate,
					PvzId:    pvz.ID,
					Status:   models.Close,
				}

				rr.On("ListByPvzIdsWithDateRange", mock.Anything, []*uuid.UUID{&pvz.ID}, startDate, endDate).
					Return([]*models.Reception{&reception}, nil)

				pr.On("ListByReceptionIds", mock.Anything, []*uuid.UUID{&reception.ID}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockPvzRepo := mocks.NewMockPvzRepository(t)
			mockReceptionRepo := mocks.NewMockReceptionRepository(t)
			mockProductRepo := mocks.NewMockProductRepository(t)

			tc.setupMocks(mockPvzRepo, mockReceptionRepo, mockProductRepo)

			pvzService := services.NewPvzService(mockPvzRepo, mockReceptionRepo, mockProductRepo)

			pvzInfo, err := pvzService.PvzInfo(context.Background(), tc.dto)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, pvzInfo)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, pvzInfo)
			}

			mockPvzRepo.AssertExpectations(t)
			mockReceptionRepo.AssertExpectations(t)
			mockProductRepo.AssertExpectations(t)
		})
	}
}
