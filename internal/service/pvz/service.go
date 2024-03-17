package pvz

import (
	"HW1/internal/model/pvz"
	pvzStorage "HW1/internal/storage/pvz"
	"fmt"
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
