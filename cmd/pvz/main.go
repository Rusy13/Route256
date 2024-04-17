package main

import (
	"log"
	"sync"

	pvzService "Homework/internal/service/pvz"
	"Homework/internal/storage/pvz"
	pvzCall "Homework/utils/call/pvz"
)

func main() {
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

	createCmdCh := make(chan []string)
	listCmdCh := make(chan struct{})

	SignCh := make(chan string)

	go stor.HandleSignals()
	serv := pvzService.New(&stor)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		cli := pvzCall.NewCLI(serv)
		cli.List(listCmdCh, SignCh)
	}()

	go func() {
		cli := pvzCall.NewCLI(serv)
		cli.Create(createCmdCh, SignCh)
	}()

	go func() {
		cli := pvzCall.NewCLI(serv)
		cli.Run(createCmdCh, listCmdCh, SignCh)
	}()

	go func() {
		defer wg.Done()
		pvzService.MonitorThreads(SignCh)
	}()
	wg.Wait()
}
