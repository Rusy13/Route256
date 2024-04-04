package order

import (
	pvzService "HW1/internal/service/pvz"
	"testing"
	"time"
)

func TestCreateCommand(t *testing.T) {
	serv := pvzService.New(nil) // Здесь может быть любая реализация StorageI, но nil используется для тестирования
	cli := NewCLI(serv)
	createCmdCh := make(chan []string)
	SignCh := make(chan string)

	go cli.Create(createCmdCh, SignCh)

	// Отправляем команду через канал
	createCmdCh <- []string{"PvzName", "Address", "Email"}

	// Проверяем, что команда успешно отправлена
	select {
	case <-SignCh:
		// Команда успешно отправлена
	case <-time.After(time.Second):
		t.Error("время ожидания истекло")
	}
}

func TestListCommand(t *testing.T) {
	serv := pvzService.New(nil) // Здесь может быть любая реализация StorageI, но nil используется для тестирования
	cli := NewCLI(serv)
	listCmdCh := make(chan struct{})
	SignCh := make(chan string)

	go cli.List(listCmdCh, SignCh)

	// Отправляем команду через канал
	listCmdCh <- struct{}{}

	// Проверяем, что команда успешно отправлена
	select {
	case <-SignCh:
		// Команда успешно отправлена
	case <-time.After(time.Second):
		t.Error("время ожидания истекло")
	}
}
