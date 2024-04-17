package main

import (
	"log"

	orderService "Homework/internal/service/orderserv"
	"Homework/internal/storage/order"
	orderCall "Homework/utils/call/order"
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
