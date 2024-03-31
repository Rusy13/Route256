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
	OrderCost   int       `json:"storageCost"`
	OrderWeight int       `json:"storageWeight"`
}

// IssueOrdersInput используется для выдачи заказов клиенту
type IssueOrdersInput struct {
	ClientID int   `json:"clientId"`
	OrderIDs []int `json:"orderIds"`
}

// PackageType определяет тип упаковки
type PackageType string

const (
	Box    PackageType = "BOX"
	Packet PackageType = "PACKET"
	Tape   PackageType = "TAPE"
)

type PackageParams struct {
	Price          int
	Name           string
	MaxOrderWeight int
}

type Packager interface {
	GetParams() PackageParams
}
