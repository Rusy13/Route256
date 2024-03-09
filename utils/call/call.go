package call

import (
	models "HW1/internal/model"
	"HW1/internal/service"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CLI struct {
	service service.Service
}

func NewCLI(serv service.Service) *CLI {
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
	case "refund":
		cli.Refund(args)
	case "issue":
		cli.Issue(args)
	case "list":
		cli.List(args)
	case "acceptrefund":
		cli.AcceptRefund(args)
	case "refundlist":
		cli.RefundList(args)
	case "help":
		cli.Help()
	default:
		fmt.Println("неизвестная команда")
	}
}

func (cli *CLI) Create(args []string) {
	// Реализация логики для команды "create"
	if len(args) < 3 {
		fmt.Println("необходимо указать ID заказа, ID получателя и срок хранения")
		return
	}
	orderID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("неверный формат ID заказа:", err)
		return
	}
	clientID, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("неверный формат ID получателя:", err)
		return
	}
	storageTimeSec, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("неверный формат срока хранения:", err)
		return
	}
	// Преобразование секунд в объект time.Duration
	duration := time.Duration(storageTimeSec) * time.Second
	// Создание временного объекта с текущим временем и добавлением продолжительности
	storageTime := time.Now().Add(duration)
	err = cli.service.AcceptOrderFromCourier(models.OrderInput{OrderID: orderID, ClientID: clientID, StorageTime: storageTime})
	if err != nil {
		fmt.Println("ошибка при принятии заказа:", err)
	} else {
		fmt.Println("заказ успешно принят")
	}
}

func (cli *CLI) Refund(args []string) {
	// Реализация логики для команды "refund"
	if len(args) < 1 {
		fmt.Println("необходимо указать ID заказа для возврата")
		return
	}
	orderID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("неверный формат ID заказа:", err)
		return
	}
	err = cli.service.ReturnOrderToCourier(orderID)
	if err != nil {
		fmt.Println("ошибка при возврате заказа:", err)
	} else {
		fmt.Println("заказ успешно возвращен курьеру")
	}
}

func (cli *CLI) Issue(args []string) {
	// Реализация логики для команды "issue"
	if len(args) < 1 {
		fmt.Println("необходимо указать ID заказов для выдачи")
		return
	}
	orderIDs := make([]int, len(args))
	for i, arg := range args {
		orderID, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("неверный формат ID заказа %s: %v\n", arg, err)
			return
		}
		orderIDs[i] = orderID
	}
	err := cli.service.IssueOrderToClient(orderIDs)
	if err != nil {
		fmt.Println("ошибка при выдаче заказов:", err)
	} else {
		fmt.Println("заказы успешно выданы клиенту")
	}
}

func (cli *CLI) List(args []string) {
	// Реализация логики для команды "list"
	if len(args) < 1 {
		fmt.Println("необходимо указать ID пользователя")
		return
	}
	clientID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("неверный формат ID пользователя:", err)
		return
	}
	var lastNOrders int
	var inPvz bool
	if len(args) > 1 {
		for _, arg := range args[1:] {
			if strings.HasPrefix(arg, "lastN=") {
				lastNOrdersStr := strings.TrimPrefix(arg, "lastN=")
				lastNOrders, err = strconv.Atoi(lastNOrdersStr)
				if err != nil {
					fmt.Println("неверный формат количества последних заказов:", err)
					return
				}
			} else if arg == "inPvz" {
				inPvz = true
			}
		}
	}
	orders, err := cli.service.GetOrderList(clientID, lastNOrders, inPvz)
	if err != nil {
		fmt.Println("ошибка при получении списка заказов:", err)
		return
	}
	fmt.Println("список заказов:")
	for _, order := range orders {
		fmt.Println(order)
	}
}

func (cli *CLI) AcceptRefund(args []string) {
	// Реализация логики для команды "acceptrefund"
	if len(args) < 2 {
		fmt.Println("необходимо указать ID пользователя и ID заказа для возврата")
		return
	}
	clientID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("неверный формат ID пользователя:", err)
		return
	}
	orderID, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("неверный формат ID заказа:", err)
		return
	}
	err = cli.service.AcceptRefundFromClient(orderID, clientID)
	if err != nil {
		fmt.Println("ошибка при принятии возврата заказа:", err)
	} else {
		fmt.Println("возврат заказа успешно принят")
	}
}

func (cli *CLI) RefundList(args []string) {
	// Реализация логики для команды "refundlist"
	if len(args) < 2 {
		fmt.Println("необходимо указать номер начального заказа и количество заказов")
		return
	}
	firstNumber, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("неверный формат номера начального заказа:", err)
		return
	}
	numberOfOrders, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("неверный формат количества заказов:", err)
		return
	}
	returnOrders, err := cli.service.GetRefundList(firstNumber, numberOfOrders)
	if err != nil {
		fmt.Println("ошибка при получении списка возвратов:", err)
		return
	}
	fmt.Println("список возвратов:")
	for _, order := range returnOrders {
		fmt.Println(order)
	}
}

func (cli *CLI) Help() {
	// Реализация вывода справки
	printHelp()
}

func printHelp() {
	fmt.Println("список доступных команд:")
	fmt.Println("create <orderID> <clientID> <storageTimeSec>: Принять заказ от курьера")
	fmt.Println("refund <orderID>: Вернуть заказ курьеру")
	fmt.Println("issue <orderID1> <orderID2> ...: Выдать заказ клиенту")
	fmt.Println("list <userID> [lastN=<N>] [inPvz]: Получить список заказов")
	fmt.Println("acceptrefund <userID> <orderID>: Принять возврат от клиента")
	fmt.Println("refundlist <pageNumber> <ordersPerPage>: Получить список возвратов")
}
