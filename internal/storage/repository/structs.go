package repository

import (
	"errors"
)

var ErrObjectNotFound = errors.New("not found")

type Pvz struct {
	ID      int64  `db:"id"`
	PvzName string `db:"pvzname"`
	Address string `db:"address"`
	Email   string `db:"email"`
}
