package pvz

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"Homework/internal/model/pvz"
)

const storageName = "storagePvz"

type Storage struct {
	storage *os.File
	pvzs    []PvzDTO
	//muRead  sync.Mutex
	//muWrite sync.Mutex
	mu sync.RWMutex
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
	//s.muRead.Lock()
	//defer s.muRead.Unlock()
	s.mu.RLock()
	defer s.mu.RUnlock()
	//time.Sleep(15 * time.Second)
	return s.pvzs, nil
}

func (s *Storage) Create(input pvz.Pvz) error {
	//s.muWrite.Lock()
	//defer s.muWrite.Unlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	//time.Sleep(15 * time.Second)

	for _, pvz := range s.pvzs {
		if pvz.PvzName == input.PvzName {
			return errors.New("пвз уже принят")
		}
	}

	// Создаем новый PVZ
	newPvz := PvzDTO{
		PvzName: input.PvzName,
		Address: input.Address,
		Email:   input.Email,
		AddDate: time.Now(),
	}

	s.pvzs = append(s.pvzs, newPvz)
	err := s.writeBytes()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) writeBytes() error {

	rawBytes, err := json.Marshal(s.pvzs)
	if err != nil {
		return err
	}

	err = os.WriteFile(storageName, rawBytes, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close() error {
	err := s.storage.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) HandleSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		sig := <-sigCh
		switch sig {
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGTERM:
			err := s.Close()
			if err != nil {
				panic(err)
			}
			time.Sleep(5 * time.Second)
			fmt.Println()
			fmt.Println("завершение прошло успешно")
			os.Exit(0)
		}
	}
}
