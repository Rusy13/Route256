package repository

import "context"

type PvzRepo interface {
	Add(ctx context.Context, pvz *Pvz) (int64, error)
	GetByID(ctx context.Context, id int64) (*Pvz, error)
	Update(ctx context.Context, id int64, pvz *Pvz) error
	Delete(ctx context.Context, id int64) error
}
