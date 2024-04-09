package pvz

import (
	"HW1/internal/model/pvz"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Errorf("Ошибка при создании нового хранилища: %v", err)
	}

	err = os.Remove(storageName)
	if err != nil {
		t.Errorf("Ошибка при удалении файла хранилища: %v", err)
	}
}

func TestStorage_ListAll(t *testing.T) {
	storage, err := New()
	if err != nil {
		t.Fatalf("Ошибка при создании нового хранилища: %v", err)
	}

	defer func() {
		err := os.Remove(storageName)
		if err != nil {
			t.Fatalf("Ошибка при удалении файла хранилища: %v", err)
		}
	}()

	expectedPvz := PvzDTO{
		PvzName: "TestPVZ",
		Address: "TestAddress",
		Email:   "test@example.com",
		AddDate: time.Now(),
	}
	storage.pvzs = append(storage.pvzs, expectedPvz)

	pvzs, err := storage.ListAll()
	if err != nil {
		t.Errorf("Ошибка при получении списка ПВЗ: %v", err)
	}

	expectedLength := 1
	if len(pvzs) != expectedLength {
		t.Errorf("Ожидалось %d ПВЗ в списке, получено: %d", expectedLength, len(pvzs))
	}

	if pvzs[0] != expectedPvz {
		t.Errorf("Ожидалось ПВЗ %+v, получено: %+v", expectedPvz, pvzs[0])
	}
}

func TestStorage_Create(t *testing.T) {
	storage, err := New()
	if err != nil {
		t.Fatalf("Ошибка при создании нового хранилища: %v", err)
	}

	defer func() {
		err := os.Remove(storageName)
		if err != nil {
			t.Fatalf("Ошибка при удалении файла хранилища: %v", err)
		}
	}()

	newPvz := pvz.Pvz{
		PvzName: "NewPVZ",
		Address: "NewAddress",
		Email:   "new@example.com",
	}

	err = storage.Create(newPvz)
	if err != nil {
		t.Errorf("Ошибка при создании ПВЗ: %v", err)
	}

	if len(storage.pvzs) != 1 {
		t.Errorf("Ожидалось добавление 1 ПВЗ в хранилище, получено: %d", len(storage.pvzs))
	}
}
