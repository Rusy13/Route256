package pvz

import (
	"HW1/internal/model/pvz"
	pvzStorage "HW1/internal/storage/pvz"
	"fmt"
	"testing"
)

type mockStorage struct{}

func (m mockStorage) Create(input pvz.Pvz) error {
	return nil
}

func (m mockStorage) ListAll() ([]pvzStorage.PvzDTO, error) {
	return nil, nil
}

func (m mockStorage) HandleSignals() {}

func TestNew(t *testing.T) {
	storage := mockStorage{}
	service := New(storage)

	if service.storage != storage {
		t.Error("Ожидается, что новый сервис будет использовать переданное хранилище")
	}
}

func TestService_CreatePvz(t *testing.T) {
	mockStorage := mockStorage{}
	service := New(mockStorage)

	err := service.CreatePvz(pvz.Pvz{})
	if err != nil {
		t.Errorf("Ожидалось успешное создание ПВЗ, получена ошибка: %v", err)
	}
	fmt.Println(service)
}

func TestService_GetPvzList(t *testing.T) {
	mockStorage := mockStorage{}
	service := New(mockStorage)

	_, err := service.GetPvzList()
	if err != nil {
		t.Errorf("Ожидалось успешное получение списка ПВЗ, получена ошибка: %v", err)
	}
}
