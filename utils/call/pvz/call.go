package order

import (
	models "HW1/internal/model/pvz"
	"HW1/internal/service/pvz"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type CLI struct {
	service pvz.Service
}

func NewCLI(serv pvz.Service) *CLI {
	return &CLI{service: serv}
}

func (cli *CLI) Run(createCmdCh chan<- []string, listCmdCh chan<- struct{}, SignCh chan<- string) {
	go func() {
		fmt.Println("Введите команду (create <PvzName> <Address> <Email> или list):")
		reader := bufio.NewReader(os.Stdin)

		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				log.Println("ошибка чтения команды:", err)
				continue
			}
			text = strings.TrimSpace(text)
			args := strings.Fields(text)
			if len(args) == 0 {
				fmt.Println("необходимо указать команду")
				continue
			}
			command := args[0]
			switch command {
			case "create":
				if len(args) < 4 {
					fmt.Println("необходимо указать название ПВЗ, адрес и email")
					continue
				}
				select {
				case createCmdCh <- args[1:]:
					SignCh <- "Команда create отправлена"
				default:
					SignCh <- "Команда create уже в обработке"
				}
			case "list":
				select {
				case listCmdCh <- struct{}{}:
					SignCh <- "Команда list отправлена"
				default:
					SignCh <- "Команда list уже в обработке"
				}
			default:
				fmt.Println("неизвестная команда")
			}
		}
	}()
}

func (cli *CLI) Create(createCmdCh <-chan []string, SignCh chan<- string) {
	for args := range createCmdCh {
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
		SignCh <- "create завершен"
	}
}

func (cli *CLI) List(listCmdCh <-chan struct{}, SignCh chan<- string) {

	for range listCmdCh {
		orders, err := cli.service.GetPvzList()
		if err != nil {
			fmt.Println("ошибка при получении списка пвз:", err)
			return
		}
		fmt.Println("список пвз:")
		for _, order := range orders {
			fmt.Println(order)
		}
		SignCh <- "list завершен"
	}
}

func (cli *CLI) Help() {
	printHelp()
}

func printHelp() {
	fmt.Println("список доступных команд:")
	fmt.Println("create <PvzName> <Address> <Email>: Записать ПВЗ")
	fmt.Println("list: Получить список ПВЗ")
}
