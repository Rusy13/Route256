package main

import (
	pvz2 "HW1/internal/service/pvz"
	pvz4 "HW1/internal/service/pvz"
	"HW1/internal/storage/pvz"
	pvz3 "HW1/utils/call/pvz"
	"log"
	"sync"
)

func main() {
	// Создаем экземпляр хранилища
	stor, err := pvz.New()
	if err != nil {
		log.Fatal("не удалось подключиться к хранилищу:", err)
	}
	defer func() {
		err := stor.Close()
		if err != nil {
			log.Println("ошибка при закрытии хранилища:", err)
		}
	}()

	go stor.HandleSignals()
	// Создаем сервис, передавая ему хранилище
	serv := pvz2.New(&stor)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		cli := pvz3.NewCLI(serv)
		cli.Run()
	}()

	go pvz4.MonitorThreads()
	wg.Wait()
}
