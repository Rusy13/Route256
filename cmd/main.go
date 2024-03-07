package main

import (
	"HW1/internal/service"
	"HW1/internal/storage"
	"HW1/utils/call"
	"fmt"
	"log"
	"os"
)

func main() {
	// Создаем экземпляр хранилища
	stor, err := storage.New()
	if err != nil {
		log.Fatal("не удалось подключиться к хранилищу:", err)
	}
	// Создаем сервис, передавая ему хранилище
	serv := service.New(&stor)

	// Проверяем наличие команды в аргументах командной строки
	if len(os.Args) < 2 {
		fmt.Println("необходимо указать команду")
		return
	}
	command := os.Args[1]
	args := os.Args[2:]

	call.CallPrograms(command, args, serv)

}
