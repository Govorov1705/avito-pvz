package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Govorov1705/avito-pvz/internal/dtos"
	"github.com/Govorov1705/avito-pvz/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntireFlow(t *testing.T) {
	// Создание нового ПВЗ
	ctx := context.Background()

	moderatorToken, err := userService.DummyLogin(ctx, &dtos.DummyLoginRequest{
		Role: "moderator",
	})
	require.NoError(t, err)

	createPvzRequest := dtos.CreatePvzRequest{
		City: "Москва",
	}
	createPvzRequestJson, err := json.Marshal(createPvzRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/pvz", bytes.NewBuffer(createPvzRequestJson))
	req.Header.Set("Authorization", "Bearer "+moderatorToken)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	pvz := models.Pvz{}
	err = json.Unmarshal(w.Body.Bytes(), &pvz)
	assert.NoError(t, err)

	// Добавление новой приемки
	employeeToken, err := userService.DummyLogin(ctx, &dtos.DummyLoginRequest{
		Role: "employee",
	})
	assert.NoError(t, err)

	createReceptionRequest := dtos.CreateReceptionRequest{PvzId: pvz.ID}
	createReceptionRequestJson, err := json.Marshal(createReceptionRequest)
	assert.NoError(t, err)

	req = httptest.NewRequest("POST", "/receptions", bytes.NewBuffer(createReceptionRequestJson))
	req.Header.Set("Authorization", "Bearer "+employeeToken)

	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	reception := models.Reception{}
	err = json.Unmarshal(w.Body.Bytes(), &reception)
	assert.NoError(t, err)

	// Добавление 50 товаров
	for i := 0; i < 50; i++ {
		addProductRequest := dtos.AddProductRequest{
			Type:  "одежда",
			PvzId: pvz.ID,
		}

		addProductRequestJson, err := json.Marshal(addProductRequest)
		assert.NoError(t, err)

		req = httptest.NewRequest("POST", "/products", bytes.NewBuffer(addProductRequestJson))
		req.Header.Set("Authorization", "Bearer "+employeeToken)

		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)

		product := models.Product{}
		err = json.Unmarshal(w.Body.Bytes(), &product)
		assert.NoError(t, err)
	}

	// Закрытие приемки
	req = httptest.NewRequest(
		"POST",
		fmt.Sprintf("/pvz/%v/close_last_reception", pvz.ID.String()),
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+employeeToken)

	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	reception = models.Reception{}
	err = json.Unmarshal(w.Body.Bytes(), &reception)
	assert.NoError(t, err)
	assert.Equal(t, models.Close, reception.Status)
}
