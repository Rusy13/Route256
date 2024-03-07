package call

import (
	models "HW1/internal/model"
	"HW1/internal/service"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func CallPrograms(command string, args []string, serv service.Service) {
	switch command {
	case "create":
		if len(args) < 3 {
			fmt.Println("Необходимо указать ID заказа, ID получателя и срок хранения")
			return
		}
		orderID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Неверный формат ID заказа:", err)
			return
		}
		clientID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Неверный формат ID получателя:", err)
			return
		}
		storageTimeSec, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Неверный формат срока хранения:", err)
			return
		}
		// Преобразование секунд в объект time.Duration
		duration := time.Duration(storageTimeSec) * time.Second
		// Создание временного объекта с текущим временем и добавлением продолжительности
		storageTime := time.Now().Add(duration)
		err = serv.AcceptOrderFromCourier(models.OrderInput{OrderID: orderID, ClientID: clientID, StorageTime: storageTime})
		if err != nil {
			fmt.Println("Ошибка при принятии заказа:", err)
		} else {
			fmt.Println("Заказ успешно принят")
		}

	case "refund":
		if len(args) < 1 {
			fmt.Println("Необходимо указать ID заказа для возврата")
			return
		}
		orderID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Неверный формат ID заказа:", err)
			return
		}
		err = serv.ReturnOrderToCourier(orderID)
		if err != nil {
			fmt.Println("Ошибка при возврате заказа:", err)
		} else {
			fmt.Println("Заказ успешно возвращен курьеру")
		}

	case "issue":
		if len(args) < 1 {
			fmt.Println("Необходимо указать ID заказов для выдачи")
			return
		}
		orderIDs := make([]int, len(args))
		for i, arg := range args {
			orderID, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Неверный формат ID заказа %s: %v\n", arg, err)
				return
			}
			orderIDs[i] = orderID
		}
		err := serv.IssueOrderToClient(orderIDs)
		if err != nil {
			fmt.Println("Ошибка при выдаче заказов:", err)
		} else {
			fmt.Println("Заказы успешно выданы клиенту")
		}

	case "list":
		if len(args) < 1 {
			fmt.Println("Необходимо указать ID пользователя")
			return
		}
		idClient, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Неверный формат ID пользователя:", err)
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
						fmt.Println("Неверный формат количества последних заказов:", err)
						return
					}
				} else if arg == "inPvz" {
					inPvz = true
				}
			}
		}
		orders, err := serv.GetOrderList(idClient, lastNOrders, inPvz)
		if err != nil {
			fmt.Println("Ошибка при получении списка заказов:", err)
			return
		}
		fmt.Println("Список заказов:")
		for _, order := range orders {
			fmt.Println(order)
		}

	case "acceptrefund":
		if len(args) < 2 {
			fmt.Println("Необходимо указать ID пользователя и ID заказа для возврата")
			return
		}
		idClient, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Неверный формат ID пользователя:", err)
			return
		}
		idOrder, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Неверный формат ID заказа:", err)
			return
		}
		err = serv.AcceptRefundFromClient(idOrder, idClient)
		if err != nil {
			fmt.Println("Ошибка при принятии возврата заказа:", err)
		} else {
			fmt.Println("Возврат заказа успешно принят")
		}

	case "refundlist":
		if len(args) < 2 {
			fmt.Println("Необходимо указать номер страницы и количество заказов на странице")
			return
		}
		firstNumber, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Неверный формат номера страницы:", err)
			return
		}
		numberOfOrders, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Неверный формат количества заказов на странице:", err)
			return
		}
		returnOrders, err := serv.GetRefundList(firstNumber, numberOfOrders)
		if err != nil {
			fmt.Println("Ошибка при получении списка возвратов:", err)
			return
		}
		fmt.Println("Список возвратов:")
		for _, order := range returnOrders {
			fmt.Println(order)
		}

	case "help":
		// Вывод списка доступных команд с описанием
		printHelp()

	default:
		fmt.Println("Неизвестная команда")
	}
}

func printHelp() {
	fmt.Println("Список доступных команд:")
	fmt.Println("create <orderID> <clientID> <storageTimeSec>: Принять заказ от курьера")
	fmt.Println("refund <orderID>: Вернуть заказ курьеру")
	fmt.Println("issue <orderID1> <orderID2> ...: Выдать заказ клиенту")
	fmt.Println("list <userID> [lastN=<N>] [inPvz]: Получить список заказов")
	fmt.Println("acceptrefund <userID> <orderID>: Принять возврат от клиента")
	fmt.Println("refundlist <pageNumber> <ordersPerPage>: Получить список возвратов")
}
