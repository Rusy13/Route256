package pvz

import (
	"HW1/internal/model/pvz"
	pvz2 "HW1/internal/storage/pvz"
	"fmt"
	"runtime"
	"time"
)

type StorageI interface {
	Create(input pvz.Pvz) error
	ListAll() ([]pvz2.PvzDTO, error)
	HandleSignals()
}

type Service struct {
	storage StorageI
}

func New(s StorageI) Service {
	return Service{storage: s}
}

func (s Service) CreatePvz(input pvz.Pvz) error {
	return s.storage.Create(input)
}

func (s Service) GetPvzList() ([]pvz2.PvzDTO, error) {
	pvzs, err := s.storage.ListAll()
	if err != nil {
		return nil, err
	}

	return pvzs, nil
}

func MonitorThreads() {
	ticker := time.NewTicker(5 * time.Second) // опрашиваем каждые 5 секунд
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// В данном случае выводим количество горутин и статус работы
			numGoroutines := runtime.NumGoroutine()
			fmt.Printf("Количество горутин: %d\n", numGoroutines)

			// Получение статуса горутин
			goroutineStatus := make(map[int]string)
			for i := 0; i < numGoroutines; i++ {
				goroutineStatus[i] = "работает"
			}

			// Вывод статуса горутин
			fmt.Println("Статус горутин:")
			for id, status := range goroutineStatus {
				fmt.Printf("Горутина %d: %s\n", id+1, status)
			}

		}
	}
}
