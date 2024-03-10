package order

import (
	models "HW1/internal/model/pvz"
	"HW1/internal/service/pvz"
	"fmt"
	"os"
)

type CLI struct {
	service pvz.Service
}

func NewCLI(serv pvz.Service) *CLI {
	return &CLI{service: serv}
}

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("необходимо указать команду")
		return
	}
	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "create":
		cli.Create(args)
	case "list":
		cli.List()
	case "help":
		cli.Help()
	default:
		fmt.Println("неизвестная команда")
	}
}

func (cli *CLI) Create(args []string) {
	// Реализация логики для команды "create"
	if len(args) < 3 {
		fmt.Println("необходимо указать название пвз, адрес и эмаил")
		return
	}
	name := args[0]
	address := args[1]
	email := args[2]

	err := cli.service.CreatePvz(models.Pvz{PvzName: name, Address: address, Email: email})
	if err != nil {
		fmt.Println("ошибка создании ПВЗ:", err)
	} else {
		fmt.Println("ПВЗ успешно создан")
	}
}

func (cli *CLI) List() {
	// Реализация логики для команды "list"
	orders, err := cli.service.GetPvzList()
	if err != nil {
		fmt.Println("ошибка при получении списка пвз:", err)
		return
	}
	fmt.Println("список пвз:")
	for _, order := range orders {
		fmt.Println(order)
	}
}

func (cli *CLI) Help() {
	// Реализация вывода справки
	printHelp()
}

func printHelp() {
	fmt.Println("список доступных команд:")
	fmt.Println("create <PvzName> <Address> <Email>: Записать ПВЗ")
	fmt.Println("list: Получить список ПВЗ")
}
