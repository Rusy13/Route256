package order

import (
	order2 "HW1/internal/model/order"
	"HW1/internal/storage/order"
	"errors"
	"fmt"
	"math"
	"time"
)

type StorageI interface {
	Create(input order2.OrderInput) error
	Delete(id int) error
	Refund(clientID int, orderID int) error
	ListAll() ([]order.OrderDTO, error)
	Issued(ordersList map[int]bool, err error) error
}

type Service struct {
	storage StorageI
}

type PackageParams struct {
	Price          int
	Name           string
	MaxOrderWeight int
}

type PackageType map[string]PackageParams

var packages = PackageType{
	"packet": {Price: 5, Name: "packet", MaxOrderWeight: 10},
	"box":    {Price: 20, Name: "box", MaxOrderWeight: 30},
	"tape":   {Price: 1, Name: "tape", MaxOrderWeight: math.MaxInt64},
}

func GetPackageParams(name string) (PackageParams, error) {
	params, ok := packages[name]
	if !ok {
		return PackageParams{}, fmt.Errorf("упаковка %s не найдена", name)
	}
	return params, nil
}

func New(s StorageI) Service {
	return Service{storage: s}
}

func (s Service) AcceptOrderFromCourier(input order2.OrderInput, packageName string) error {
	if input.StorageTime.Before(time.Now()) {
		return errors.New("срок хранения в прошлом")
	}
	packageParams, err := GetPackageParams(packageName)

	if packageParams.MaxOrderWeight < input.OrderWeight {
		return fmt.Errorf("вес заказа для упаковки не подходит")
	}

	if err != nil {
		return fmt.Errorf("ошибка при получении параметров упаковки: %v", err)
	}

	input.OrderCost += packageParams.Price
	return s.storage.Create(input)
}

func (s Service) ReturnOrderToCourier(orderID int) error {
	orders, err := s.storage.ListAll()
	if err != nil {
		return err
	}
	if orderID <= 0 {
		return errors.New("неверный формат id заказа (id может быть только больше 0)")
	}
	for _, order := range orders {
		if (order.OrderID == orderID) && (order.StorageTime.Before(time.Now())) && (!order.IsIssued) {
			//orders[index].IsDeleted = true
			return s.storage.Delete(orderID)
		}
	}
	return errors.New("срок действия не истек") //s.s.Delete(orderID)
}

func (s Service) IssueOrderToClient(orderIDs []int) error {
	orders, err := s.storage.ListAll()
	fmt.Println(orders)
	if err != nil {
		return err
	}
	ordersMap := make(map[int]bool)
	for _, order := range orderIDs {
		ordersMap[order] = false
	}

	prevIdClient := 0
	for _, allorder := range orders {
		_, ok := ordersMap[allorder.OrderID]
		if ok && time.Now().Before(allorder.StorageTime) {
			if prevIdClient == 0 {
				ordersMap[allorder.OrderID] = true
				prevIdClient = allorder.ClientID
			} else if allorder.ClientID == prevIdClient {
				ordersMap[allorder.OrderID] = true
			} else {
				return errors.New("у заказа разные получатели")
			}
		} else {
			return errors.New("заказ не найден")
		}

	}
	return s.storage.Issued(ordersMap, err)
}

func (s Service) GetOrderList(clientID int, optionalParams ...interface{}) ([]order.OrderDTO, error) {
	orders, err := s.storage.ListAll()
	if err != nil {
		return nil, err
	}

	var lastNOrders int
	var inPvz bool

	for _, param := range optionalParams {
		switch p := param.(type) {
		case int:
			lastNOrders = p
		case bool:
			inPvz = p
		default:
			return nil, errors.New("недопустимый тип опционального параметра")
		}
	}

	if lastNOrders > 0 && lastNOrders < len(orders) {
		ordersList := make([]order.OrderDTO, 0)
		count := 0
		for _, order := range orders {
			if order.ClientID == clientID {
				ordersList = append(ordersList, order)
				count += 1
			}
		}
		if count < lastNOrders {
			return nil, errors.New("заказов в ПВЗ меньше вашего числа")

		}
		return ordersList[:lastNOrders], nil
	}

	if inPvz {
		var pvzOrders []order.OrderDTO
		for _, order := range orders {
			if order.MetkaPVZ == "PVZ_UGAROV_RUSLAN" && order.ClientID == clientID {
				pvzOrders = append(pvzOrders, order)
			}
		}
		return pvzOrders, nil
	}
	return orders, nil
}

func (s Service) AcceptRefundFromClient(orderID int, clientID int) error {
	if orderID <= 0 || clientID <= 0 {
		return errors.New("неверный формат id (id не может быть меньше или равен 0)")
	}
	orders, err := s.storage.ListAll()
	if err != nil {
		return err
	}
	for _, order := range orders {
		if order.OrderID == orderID && order.ClientID == clientID && order.MetkaPVZ == "PVZ_UGAROV_RUSLAN" && time.Since(order.IssuedDate) <= 2*24*time.Hour { //2 days{
			return s.storage.Refund(orderID, clientID)
		}
	}
	return errors.New("прошло слишком много времени или товар выдавался не нашим ПВЗ")
}

func (s Service) GetRefundList(firstNumber int, numberOfOrders int) ([]order.OrderDTO, error) {
	all, err := s.storage.ListAll()
	if err != nil {
		return nil, err
	}
	ordersList := make([]order.OrderDTO, 0)
	for _, order := range all {
		if order.IsReturned {
			ordersList = append(ordersList, order)
		}
	}
	if len(ordersList) < numberOfOrders {
		return nil, errors.New("количество возвращенных заказов меньше, чем вы ввели")
	}
	if numberOfOrders > len(ordersList[firstNumber-1:firstNumber+numberOfOrders-1]) {
		return nil, errors.New("количество ожидаемых заказов больше чем реальное количество на ПВЗ")
	}

	return ordersList[firstNumber-1 : firstNumber+numberOfOrders-1], nil
}
