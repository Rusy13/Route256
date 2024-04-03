package postgresql

import (
	"HW1/internal/storage/repository"
	mock "HW1/internal/storage/repository/mocks"
	"HW1/tests/fixtures"
	"bytes"
	"context"
	"encoding/json"
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
		assert.Equal(t, "{\"ID\":50001,\"PvzName\":\"asd\",\"Address\":\"asd\",\"Email\":\"asd\"}", string(result))
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

func TestCreate(t *testing.T) {
	t.Parallel()

	// Инициализация контроллера mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мока для вашего репозитория
	mockRepo := mock.NewMockPvzRepo(ctrl)

	// Создание экземпляра сервера
	server := Server1{Repo: mockRepo}

	// Подготовка входных данных
	pvzRepo := &repository.Pvz{
		PvzName: "TestPvz",
		Address: "TestAddress",
		Email:   "test@example.com",
	}

	// Ожидаемый результат
	expectedID := int64(12345)
	expectedJSON := []byte(`{"id":12345,"pvzname":"TestPvz","address":"TestAddress","email":"test@example.com"}`)
	expectedStatus := http.StatusOK

	// Настройка мока
	mockRepo.EXPECT().Add(gomock.Any(), pvzRepo).Return(expectedID, nil)

	// Выполнение функции
	pvzJson, status := server.create(context.Background(), pvzRepo)

	// Проверка результата
	assert.Equal(t, expectedStatus, status)
	assert.Equal(t, expectedJSON, pvzJson)
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
