package order

import "time"

type OrderDTO struct {
	OrderID     int
	ClientID    int
	StorageTime time.Time
	IsIssued    bool //выдать клиенту
	IsReturned  bool //возврат от клиента
	IsDeleted   bool //возврат курьеру
	MetkaPVZ    string
	IssuedDate  time.Time
}
