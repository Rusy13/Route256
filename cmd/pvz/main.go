package main

import (
	pvz2 "HW1/internal/service/pvz"
	"HW1/internal/storage/pvz"
	pvz3 "HW1/utils/call/pvz"
	"log"
)

func main() {
	// Создаем экземпляр хранилища
	stor, err := pvz.New()
	if err != nil {
		log.Fatal("не удалось подключиться к хранилищу:", err)
	}
	// Создаем сервис, передавая ему хранилище
	serv := pvz2.New(&stor)

	cli := pvz3.NewCLI(serv)
	cli.Run()
}
