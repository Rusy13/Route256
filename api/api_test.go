package postgresql

import (
	"HW1/internal/storage/repository"
	mock "HW1/internal/storage/repository/mocks"
	"HW1/tests/fixtures"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetByID(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPvz.EXPECT().GetByID(gomock.Any(), id).Return(fixtures.Pvz().Valid().P(), nil)
		result, status := s.srv.get(ctx, id)

		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"PvzName\":\"ExamplePvz1\",\"Address\":\"ExampleAddress\",\"Email\":\"example@example.com\"}", string(result))
	})
}

func Test_validateGetByID(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		result := validateGetByID(1)
		assert.True(t, result)
	})
	t.Run("fail", func(t *testing.T) {
		result := validateGetByID(-1)
		assert.False(t, result)
	})
}

func TestUpdatePvz(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPvzRepo(ctrl)
	server := Server1{
		Repo: mockRepo,
	}
	router := CreateRouter(server)

	mockRepo.EXPECT().Update(gomock.Any(), int64(1), gomock.Any()).Return(nil)
	updateData := map[string]string{
		"pvzname": "qq",
		"address": "Nevsky prospect. St.Petersburg",
		"email":   "www@spb.ru",
	}
	updateJSON, _ := json.Marshal(updateData)
	req, err := http.NewRequest("PUT", "/pvz/1", bytes.NewBuffer(updateJSON))
	if err != nil {
		t.Fatal(err)
	}

	req.SetBasicAuth("rus", "1234")
	rr := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"ID":0,"PvzName":"qq","Address":"Nevsky prospect. St.Petersburg","Email":"www@spb.ru"}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeletePvz(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPvzRepo(ctrl)
	server := Server1{
		Repo: mockRepo,
	}
	router := CreateRouter(server)

	// Устанавливаем ожидание вызова метода Delete у мок-репозитория
	mockRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)
	// Создаем запрос DELETE с установленным параметром ключа в URL
	req, err := http.NewRequest("DELETE", "/pvz/1", bytes.NewBuffer([]byte{}))

	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("rus", "1234")

	// Создаем тестовый Recorder для записи ответа
	rr := httptest.NewRecorder()

	// Маршрутизация запроса
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rr, req)

	// Проверяем статус код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	expected := "Successfully deleted"
	assert.Equal(t, expected, rr.Body.String())
}

// -----------------------------------------------------------------

// -----------------------------------------------------------------

func TestCreatePvzHandler(t *testing.T) {
	t.Parallel()

	// Подготовка тестовых данных
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := mock.NewMockPvzRepo(mockCtrl)

	srv := &Server1{
		Repo: mockRepo,
	}

	handler := http.HandlerFunc(srv.CreatePvz)

	testCases := []struct {
		name                 string
		requestBody          []byte
		mockRepoExpect       func()
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Valid request",
			requestBody: []byte(`{"PvzName":"PVZ 1","Address":"123 Main Street","Email":"pvz1@example.com"}`),
			mockRepoExpect: func() {
				mockRepo.EXPECT().Add(gomock.Any(), &repository.Pvz{
					PvzName: "PVZ 1",
					Address: "123 Main Street",
					Email:   "pvz1@example.com",
				}).Return(int64(1), nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"pvzname":"PVZ 1","address":"123 Main Street","email":"pvz1@example.com"}`,
		},
		{
			name:                 "Failed to unmarshal JSON",
			requestBody:          []byte(`invalid json`),
			mockRepoExpect:       func() {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "Failed to unmarshal JSON\n",
		},
		{
			name:        "Failed to add pvz",
			requestBody: []byte(`{"PvzName":"PVZ 1","Address":"123 Main Street","Email":"pvz1@example.com"}`),
			mockRepoExpect: func() {
				mockRepo.EXPECT().Add(gomock.Any(), &repository.Pvz{
					PvzName: "PVZ 1",
					Address: "123 Main Street",
					Email:   "pvz1@example.com",
				}).Return(int64(0), errors.New("internal server error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "Failed to add pvz\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(tc.requestBody))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			if tc.mockRepoExpect != nil {
				tc.mockRepoExpect()
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatusCode)
			}

			if rr.Body.String() != tc.expectedResponseBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expectedResponseBody)
			}
		})
	}
}
