package main

import (
	"HW1/internal/service"
	"HW1/internal/storage"
	"HW1/utils/call"
	"log"
)

func main() {
	// Создаем экземпляр хранилища
	stor, err := storage.New()
	if err != nil {
		log.Fatal("не удалось подключиться к хранилищу:", err)
	}
	// Создаем сервис, передавая ему хранилище
	serv := service.New(&stor)

	cli := call.NewCLI(serv)
	cli.Run()
}
