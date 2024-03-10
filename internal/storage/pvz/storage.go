package pvz

import (
	"HW1/internal/model/pvz"
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"
)

const storageName = "storagePvz"

type Storage struct {
	storage *os.File
	pvzs    []PvzDTO
}

func New() (Storage, error) {
	file, err := os.OpenFile(storageName, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}

	storage := Storage{
		storage: file,
	}

	storage.pvzs, err = storage.readOrdersFromStorage()
	if err != nil {
		return Storage{}, err
	}

	return storage, nil
}

func (s *Storage) readOrdersFromStorage() ([]PvzDTO, error) {
	reader := bufio.NewReader(s.storage)
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var orders []PvzDTO
	if len(rawBytes) == 0 {
		return orders, nil
	}

	err = json.Unmarshal(rawBytes, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Storage) ListAll() ([]PvzDTO, error) {
	return s.pvzs, nil // Просто возвращаем данные, уже прочитанные из файла
}

// Create creates order
func (s *Storage) Create(input pvz.Pvz) error {
	all, err := s.ListAll()
	if err != nil {
		return err
	}

	for _, pvz := range all {
		if pvz.PvzName == input.PvzName {
			return errors.New("пвз уже принят")
		}
	}

	newPvz := PvzDTO{
		PvzName: input.PvzName,
		Address: input.Address,
		Email:   input.Email,
		AddDate: time.Now(),
	}

	all = append(all, newPvz)
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

func writeBytes(Orders []PvzDTO) error {
	rawBytes, err := json.Marshal(Orders)
	if err != nil {
		return err
	}

	err = os.WriteFile(storageName, rawBytes, 0777)
	if err != nil {
		return err
	}
	return nil
}
