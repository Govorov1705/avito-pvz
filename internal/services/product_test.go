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

func TestProductService_Add(t *testing.T) {
	type testCase struct {
		name       string
		dto        *dtos.AddProductRequest
		setupMocks func(
			pr *mocks.MockProductRepository,
			rr *mocks.MockReceptionRepository,
			mdb pgxmock.PgxPoolIface,
		)
		expectError bool
	}

	testPvzId := uuid.New()
	testProductType := models.Clothes

	tcs := []testCase{
		{
			name: "OK",
			dto: &dtos.AddProductRequest{
				Type:  testProductType,
				PvzId: testPvzId,
			},
			setupMocks: func(
				pr *mocks.MockProductRepository,
				rr *mocks.MockReceptionRepository,
				mdb pgxmock.PgxPoolIface,
			) {
				mdb.ExpectBegin()

				reception := models.Reception{
					ID:       uuid.New(),
					DateTime: time.Now(),
					PvzId:    testPvzId,
					Status:   models.InProgress,
				}

				rr.On("GetOpenByPvzId", mock.Anything, testPvzId).
					Return(&reception, nil)

				pr.On("Add", mock.Anything, testProductType, reception.ID).
					Return(&models.Product{
						ID:          uuid.New(),
						DateTime:    time.Now(),
						Type:        testProductType,
						ReceptionId: reception.ID,
					}, nil)

				mdb.ExpectCommit()
			},
			expectError: false,
		},

		{
			name: "ReceptionRepository GetOpenByPvzId error",
			dto: &dtos.AddProductRequest{
				Type:  testProductType,
				PvzId: testPvzId,
			},
			setupMocks: func(
				pr *mocks.MockProductRepository,
				rr *mocks.MockReceptionRepository,
				mdb pgxmock.PgxPoolIface,
			) {
				mdb.ExpectBegin()

				rr.On("GetOpenByPvzId", mock.Anything, testPvzId).
					Return(nil, errors.New("some error"))

				mdb.ExpectRollback()
			},
			expectError: true,
		},

		{
			name: "ProductRepository Add error",
			dto: &dtos.AddProductRequest{
				Type:  testProductType,
				PvzId: testPvzId,
			},
			setupMocks: func(
				pr *mocks.MockProductRepository,
				rr *mocks.MockReceptionRepository,
				mdb pgxmock.PgxPoolIface,
			) {
				mdb.ExpectBegin()

				reception := models.Reception{
					ID:       uuid.New(),
					DateTime: time.Now(),
					PvzId:    testPvzId,
					Status:   models.InProgress,
				}

				rr.On("GetOpenByPvzId", mock.Anything, testPvzId).
					Return(&reception, nil)

				pr.On("Add", mock.Anything, testProductType, reception.ID).
					Return(nil, errors.New("some error"))

				mdb.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockProductRepo := mocks.NewMockProductRepository(t)
			mockReceptionRepo := mocks.NewMockReceptionRepository(t)
			mdb, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mdb.Close()

			tc.setupMocks(mockProductRepo, mockReceptionRepo, mdb)

			productService := services.NewProductService(
				mockProductRepo,
				mockReceptionRepo,
				mdb,
			)

			product, err := productService.Add(context.Background(), tc.dto)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, product)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, product)
			}

			mockProductRepo.AssertExpectations(t)
			mockReceptionRepo.AssertExpectations(t)
			if err := mdb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProductService_Delete(t *testing.T) {
	type testCase struct {
		name       string
		setupMocks func(
			pr *mocks.MockProductRepository,
			rr *mocks.MockReceptionRepository,
			mdb pgxmock.PgxPoolIface,
		)
		expectError bool
	}

	testPvzId := uuid.New()

	tcs := []testCase{
		{
			name: "OK",
			setupMocks: func(
				pr *mocks.MockProductRepository,
				rr *mocks.MockReceptionRepository,
				mdb pgxmock.PgxPoolIface,
			) {
				mdb.ExpectBegin()

				reception := models.Reception{
					ID:       uuid.New(),
					DateTime: time.Now(),
					PvzId:    testPvzId,
					Status:   models.InProgress,
				}

				rr.On("GetOpenByPvzId", mock.Anything, testPvzId).
					Return(&reception, nil)

				product := models.Product{
					ID:          uuid.New(),
					DateTime:    time.Now(),
					Type:        models.Clothes,
					ReceptionId: reception.ID,
				}

				pr.On("GetLastByReceptionId", mock.Anything, reception.ID).
					Return(&product, nil)

				pr.On("Delete", mock.Anything, product.ID).Return(nil)

				mdb.ExpectCommit()
			},
			expectError: false,
		},

		{
			name: "Begin Tx error",
			setupMocks: func(
				pr *mocks.MockProductRepository,
				rr *mocks.MockReceptionRepository,
				mdb pgxmock.PgxPoolIface,
			) {
				mdb.ExpectBegin().WillReturnError(errors.New("some error"))
			},
			expectError: true,
		},

		{
			name: "ReceptionRepository GetOpenByPvzId error",
			setupMocks: func(
				pr *mocks.MockProductRepository,
				rr *mocks.MockReceptionRepository,
				mdb pgxmock.PgxPoolIface,
			) {
				mdb.ExpectBegin()

				rr.On("GetOpenByPvzId", mock.Anything, testPvzId).
					Return(nil, errors.New("some error"))

				mdb.ExpectRollback()
			},
			expectError: true,
		},

		{
			name: "ProductRepository GetLastByReceptionId error",
			setupMocks: func(
				pr *mocks.MockProductRepository,
				rr *mocks.MockReceptionRepository,
				mdb pgxmock.PgxPoolIface,
			) {
				mdb.ExpectBegin()

				reception := models.Reception{
					ID:       uuid.New(),
					DateTime: time.Now(),
					PvzId:    testPvzId,
					Status:   models.InProgress,
				}

				rr.On("GetOpenByPvzId", mock.Anything, testPvzId).
					Return(&reception, nil)

				pr.On("GetLastByReceptionId", mock.Anything, reception.ID).
					Return(nil, errors.New("some error"))

				mdb.ExpectRollback()
			},
			expectError: true,
		},

		{
			name: "ProductRepository Delete error",
			setupMocks: func(
				pr *mocks.MockProductRepository,
				rr *mocks.MockReceptionRepository,
				mdb pgxmock.PgxPoolIface,
			) {
				mdb.ExpectBegin()

				reception := models.Reception{
					ID:       uuid.New(),
					DateTime: time.Now(),
					PvzId:    testPvzId,
					Status:   models.InProgress,
				}

				rr.On("GetOpenByPvzId", mock.Anything, testPvzId).
					Return(&reception, nil)

				product := models.Product{
					ID:          uuid.New(),
					DateTime:    time.Now(),
					Type:        models.Clothes,
					ReceptionId: reception.ID,
				}

				pr.On("GetLastByReceptionId", mock.Anything, reception.ID).
					Return(&product, nil)

				pr.On("Delete", mock.Anything, product.ID).Return(errors.New("some error"))

				mdb.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockProductRepo := mocks.NewMockProductRepository(t)
			mockReceptionRepo := mocks.NewMockReceptionRepository(t)
			mdb, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mdb.Close()

			tc.setupMocks(mockProductRepo, mockReceptionRepo, mdb)

			productService := services.NewProductService(
				mockProductRepo,
				mockReceptionRepo,
				mdb,
			)

			err = productService.Delete(context.Background(), testPvzId)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockProductRepo.AssertExpectations(t)
			mockReceptionRepo.AssertExpectations(t)
			if err := mdb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
