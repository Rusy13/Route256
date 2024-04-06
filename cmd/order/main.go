package main

import (
	orderService "HW1/internal/service/orderserv"
	"HW1/internal/storage/order"
	orderCall "HW1/utils/call/order"
	"log"
)

func main() {
	// Создаем экземпляр хранилища
	stor, err := order.New()
	if err != nil {
		log.Fatal("не удалось подключиться к хранилищу:", err)
	}
	// Создаем сервис, передавая ему хранилище
	serv := orderService.New(&stor)

	cli := orderCall.NewCLI(serv)
	cli.Run()
}
