package pvz

import (
	"HW1/internal/model/pvz"
	pvz2 "HW1/internal/storage/pvz"
)

type StorageI interface {
	Create(input pvz.Pvz) error
	ListAll() ([]pvz2.PvzDTO, error)
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
