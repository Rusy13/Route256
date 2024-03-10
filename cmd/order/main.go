package main

import (
	order2 "HW1/internal/service/order"
	"HW1/internal/storage/order"
	order3 "HW1/utils/call/order"
	"log"
)

func main() {
	// Создаем экземпляр хранилища
	stor, err := order.New()
	if err != nil {
		log.Fatal("не удалось подключиться к хранилищу:", err)
	}
	// Создаем сервис, передавая ему хранилище
	serv := order2.New(&stor)

	cli := order3.NewCLI(serv)
	cli.Run()
}
