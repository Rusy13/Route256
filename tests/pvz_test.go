//go:build integration
// +build integration

package tests

import (
	api "HW1/api"
	"HW1/internal/config"
	dbN "HW1/internal/storage/db"
	"HW1/internal/storage/repository"
	dbrepo "HW1/internal/storage/repository/postgresql"
	"HW1/tests/fixtures"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetPvzByID(t *testing.T) {
	var (
		room = fixtures.Pvz().Valid().P()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPvzRepo(db.DB)
		resp := 1

		var pvz = &repository.Pvz{
			ID:      1,                     // Пример значения для ID
			PvzName: "ExamplePvz1",         // Пример значения для PvzName
			Address: "ExampleAddress",      // Пример значения для Address
			Email:   "example@example.com", // Пример значения для Email
		}

		_, err := repo.Add(context.Background(), pvz)

		//require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetByID(context.Background(), int64(resp))

		//assert
		require.NoError(t, err)
		assert.Equal(t, room.PvzName, getRoom.PvzName)
		assert.Equal(t, room.Email, getRoom.Email)
		assert.Equal(t, room.Address, getRoom.Address)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPvzRepo(db.DB)
		respFail := 2

		var pvz = &repository.Pvz{
			ID:      1,                     // Пример значения для ID
			PvzName: "ExamplePvz1",         // Пример значения для PvzName
			Address: "ExampleAddress",      // Пример значения для Address
			Email:   "example@example.com", // Пример значения для Email
		}

		_, err := repo.Add(context.Background(), pvz)

		assert.NotZero(t, respFail)

		//act
		getRoom, err := repo.GetByID(context.Background(), int64(respFail))

		//assert
		require.EqualError(t, err, "not found")
		assert.Nil(t, getRoom)
	})
}

func TestCreatePvzHandler(t *testing.T) {
	// Настройка временной тестовой конфигурации для базы данных
	tempConfig := config.StorageConfig{
		Host:     "localhost",
		Port:     5432, // Порт вашей тестовой базы данных
		Username: "postgres",
		Password: "1111",
		Database: "TestRoute",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создание нового соединения с тестовой базой данных
	tempDatabase, err := dbN.NewDb(ctx, tempConfig)
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}
	defer tempDatabase.GetPool(ctx).Close()

	// Создание объекта репозитория для тестовой базы данных
	pvzRepo := dbrepo.NewPvzRepo(tempDatabase)

	// Создание объекта сервера API с использованием тестового репозитория
	server := api.Server1{Repo: pvzRepo}

	// Подготовка данных запроса
	requestBody, err := json.Marshal(map[string]string{
		"pvzname": "PVZ 1",
		"address": "123 Main Street",
		"email":   "pvz1@example.com",
	})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Создание HTTP запроса
	req, err := http.NewRequest("POST", "http://localhost:9000/pvz", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создание HTTP тестового сервера
	rr := httptest.NewRecorder()

	// Обработка запроса сервером API
	handler := http.HandlerFunc(server.CreatePvz)
	handler.ServeHTTP(rr, req)

	// Проверка кода состояния ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Чтение и анализ ответа
	respBody := rr.Body.String()
	fmt.Println("Response:", respBody)
	// Здесь вы можете выполнить дополнительные проверки, основанные на ожидаемом ответе
}
