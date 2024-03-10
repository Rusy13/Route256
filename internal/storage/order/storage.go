package order

import (
	"HW1/internal/model/order"
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"
)

const storageName = "storage"

type Storage struct {
	storage *os.File
	orders  []OrderDTO
}

func New() (Storage, error) {
	file, err := os.OpenFile(storageName, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}

	storage := Storage{
		storage: file,
	}

	storage.orders, err = storage.readOrdersFromStorage()
	if err != nil {
		return Storage{}, err
	}

	return storage, nil
}

func (s *Storage) readOrdersFromStorage() ([]OrderDTO, error) {
	reader := bufio.NewReader(s.storage)
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var orders []OrderDTO
	if len(rawBytes) == 0 {
		return orders, nil
	}

	err = json.Unmarshal(rawBytes, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Storage) ListAll() ([]OrderDTO, error) {
	return s.orders, nil // Просто возвращаем данные, уже прочитанные из файла
}

// Create creates order
func (s *Storage) Create(input order.OrderInput) error {
	all, err := s.ListAll()
	if err != nil {
		return err
	}

	for _, order := range all {
		if order.OrderID == input.OrderID {
			return errors.New("заказ уже принят")
		}
	}

	newOrder := OrderDTO{
		OrderID:     input.OrderID,
		ClientID:    input.ClientID,
		StorageTime: input.StorageTime,
		IsIssued:    false, //выдан клиенту
		IsReturned:  false, //возвращен
		MetkaPVZ:    "PVZ_UGAROV_RUSLAN",
	}

	all = append(all, newOrder)
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

func writeBytes(Orders []OrderDTO) error {
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

// Delete order
func (s *Storage) Delete(id int) error {
	all, err := s.ListAll()
	if err != nil {
		return err
	}
	for indx, order := range all {
		if order.OrderID == id {
			all[indx].IsDeleted = true
		}
	}
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil

}

// Refund order
func (s *Storage) Refund(clientID int, orderID int) error {
	all, err := s.ListAll()
	if err != nil {
		return err
	}
	for indx, order := range all {
		if order.OrderID == orderID && order.ClientID == clientID {
			all[indx].IsReturned = true
		}
	}
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Issued(ordersList map[int]bool, err error) error {
	all, err := s.ListAll()
	if err != nil {
		return err
	}
	for indx, orderAllList := range all {
		_, ok := ordersList[orderAllList.OrderID]
		if ok {
			all[indx].IsIssued = true
			all[indx].IssuedDate = time.Now()
		}
	}
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}
