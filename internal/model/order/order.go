package order

import "time"

// Order описывает структуру заказа
type Order struct {
	OrderID     int       `json:"orderId"`
	ClientID    int       `json:"clientId"`
	OrderName   string    `json:"orderName"`
	StorageTime time.Time `json:"storageTime"`
}

// OrderInput представляет входные данные для создания заказа
type OrderInput struct {
	OrderID     int       `json:"orderId"`
	ClientID    int       `json:"clientId"`
	StorageTime time.Time `json:"storageTime"`
}

// IssueOrdersInput используется для выдачи заказов клиенту
type IssueOrdersInput struct {
	ClientID int   `json:"clientId"`
	OrderIDs []int `json:"orderIds"`
}
