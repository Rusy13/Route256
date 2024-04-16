package pvz

import (
	"fmt"

	"Homework/internal/model/pvz"
	pvzStorage "Homework/internal/storage/pvz"
)

type StorageI interface {
	Create(input pvz.Pvz) error
	ListAll() ([]pvzStorage.PvzDTO, error)
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

func (s Service) GetPvzList() ([]pvzStorage.PvzDTO, error) {
	pvzs, err := s.storage.ListAll()
	if err != nil {
		return nil, err
	}

	return pvzs, nil
}

func MonitorThreads(SignCh <-chan string) {
	for msg := range SignCh {
		fmt.Println(msg)
	}
}
