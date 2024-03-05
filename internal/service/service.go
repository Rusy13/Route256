package service

import (
	"HW1/internal/model"
	"HW1/internal/storage"
	"errors"
	"fmt"
	"time"
)

type Service struct {
	s storage.StorageI
}

func New(s storage.StorageI) Service {
	return Service{s: s}
}

func (s Service) AcceptOrderFromCourier(input models.OrderInput) error {
	if input.StorageTime.Before(time.Now()) {
		return errors.New("срок хранения в прошлом")
	}
	return s.s.Create(input)
}

func (s Service) ReturnOrderToCourier(idOrder int) error {
	orders, err := s.s.ListAll()
	if err != nil {
		return err
	}
	if idOrder <= 0 {
		return errors.New("неверный формат id заказа (id может быть только больше 0)")
	}
	for _, order := range orders {
		if (order.OrderID == idOrder) && (order.StorageTime.Before(time.Now())) && (!order.IsIssued) {
			//orders[index].IsDeleted = true
			return s.s.Delete(idOrder)
		}
	}
	return errors.New("Срок действия не истек") //s.s.Delete(idOrder)
}

func (s Service) IssueOrderToClient(orderIDs []int) error {
	orders, err := s.s.ListAll()
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
	return s.s.Issued(ordersMap, err)
}

func (s Service) GetOrderList(idClient int, optionalParams ...interface{}) ([]storage.OrderDTO, error) {
	orders, err := s.s.ListAll()
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
			return nil, errors.New("Недопустимый тип опционального параметра")
		}
	}

	if lastNOrders > 0 && lastNOrders < len(orders) {
		ordersList := make([]storage.OrderDTO, 0)
		count := 0
		for _, order := range orders {
			if order.ClientID == idClient {
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
		pvzOrders := []storage.OrderDTO{}
		for _, order := range orders {
			if order.MetkaPVZ == "PVZ_UGAROV_RUSLAN" && order.ClientID == idClient {
				pvzOrders = append(pvzOrders, order)
			}
		}
		return pvzOrders, nil
	}
	return orders, nil
}

func (s Service) AcceptRefundFromClient(idOrder int, idClient int) error {
	if idOrder <= 0 || idClient <= 0 {
		return errors.New("неверный формат id (id не может быть меньше или равен 0)")
	}
	orders, err := s.s.ListAll()
	if err != nil {
		return err
	}
	for _, order := range orders {
		if order.OrderID == idOrder && order.ClientID == idClient && order.MetkaPVZ == "PVZ_UGAROV_RUSLAN" && time.Since(order.IssuedDate) <= 2*24*time.Hour { //2 days{
			return s.s.Refund(idOrder, idClient)
		}
	}
	return errors.New("Прошло слишком много времени или товар выдавался не нашим ПВЗ")
}

func (s Service) GetRefundList(firstNumber int, numberOfOrders int) ([]storage.OrderDTO, error) {
	all, err := s.s.ListAll()
	if err != nil {
		return nil, err
	}
	ordersList := make([]storage.OrderDTO, 0)
	for _, order := range all {
		if order.IsReturned {
			ordersList = append(ordersList, order)
		}
	}
	if numberOfOrders > len(ordersList[firstNumber-1:firstNumber+numberOfOrders-1]) {
		return nil, errors.New("количество ожидаемых заказов больше чем реальное количество на ПВЗ")
	}
	return ordersList[firstNumber-1 : firstNumber+numberOfOrders-1], nil
}
