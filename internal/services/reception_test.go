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
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReceptionService_CreateReception(t *testing.T) {
	type testCase struct {
		name        string
		dto         *dtos.CreateReceptionRequest
		setupMock   func(r *mocks.MockReceptionRepository)
		expectError bool
	}

	testPvzId := uuid.New()

	tcs := []testCase{
		{
			name: "OK",
			dto: &dtos.CreateReceptionRequest{
				PvzId: testPvzId,
			},
			setupMock: func(r *mocks.MockReceptionRepository) {
				r.On("Add", mock.Anything, &dtos.CreateReceptionRequest{
					PvzId: testPvzId,
				}).Return(&models.Reception{
					ID:       uuid.New(),
					DateTime: time.Now(),
					PvzId:    testPvzId,
					Status:   models.InProgress,
				}, nil)
			},
			expectError: false,
		},

		{
			name: "ReceptionRepository Add error",
			dto: &dtos.CreateReceptionRequest{
				PvzId: testPvzId,
			},
			setupMock: func(r *mocks.MockReceptionRepository) {
				r.On("Add", mock.Anything, &dtos.CreateReceptionRequest{
					PvzId: testPvzId,
				}).Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockReceptionRepo := mocks.NewMockReceptionRepository(t)

			tc.setupMock(mockReceptionRepo)

			receptionService := services.NewReceptionService(mockReceptionRepo, nil)

			reception, err := receptionService.CreateReception(context.Background(), tc.dto)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, reception)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, reception)
			}

			mockReceptionRepo.AssertExpectations(t)
		})
	}
}

func TestReceptionService_CloseReception(t *testing.T) {
	type testCase struct {
		name       string
		setupMocks func(
			mdb pgxmock.PgxPoolIface,
			rr *mocks.MockReceptionRepository,
		)
		expectError bool
	}

	testPvzId := uuid.New()
	testReceptionId := uuid.New()
	testReceptionDateTime := time.Now()

	tcs := []testCase{
		{
			name: "OK",
			setupMocks: func(
				mdb pgxmock.PgxPoolIface,
				rr *mocks.MockReceptionRepository,
			) {
				mdb.ExpectBegin()

				rr.On("GetOpenByPvzId", mock.Anything, testPvzId).Return(&models.Reception{
					ID:       testReceptionId,
					DateTime: testReceptionDateTime,
					PvzId:    testPvzId,
					Status:   models.InProgress,
				}, nil)

				rr.On("Close", mock.Anything, testReceptionId).
					Return(&models.Reception{
						ID:       testReceptionId,
						DateTime: testReceptionDateTime,
						PvzId:    testPvzId,
						Status:   models.Close,
					}, nil)

				mdb.ExpectCommit()
			},
			expectError: false,
		},

		{
			name: "Error starting transaction",
			setupMocks: func(
				mdb pgxmock.PgxPoolIface,
				rr *mocks.MockReceptionRepository,
			) {
				mdb.ExpectBegin().WillReturnError(errors.New("some error"))
			},
			expectError: true,
		},

		{
			name: "ReceptionRepository GetOpenByPvzId error",
			setupMocks: func(
				mdb pgxmock.PgxPoolIface,
				rr *mocks.MockReceptionRepository,
			) {
				mdb.ExpectBegin()

				rr.On("GetOpenByPvzId", mock.Anything, testPvzId).Return(nil, errors.New("some error"))

				mdb.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockReceptionRepo := mocks.NewMockReceptionRepository(t)
			mdb, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mdb.Close()

			tc.setupMocks(mdb, mockReceptionRepo)

			receptionService := services.NewReceptionService(mockReceptionRepo, mdb)

			reception, err := receptionService.CloseReception(context.Background(), testPvzId)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, reception)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, reception)
			}

			mockReceptionRepo.AssertExpectations(t)
			if err := mdb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
