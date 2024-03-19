package postgresql

import (
	"HW1/pkg/db"
	"HW1/pkg/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PvzRepo struct {
	db *db.Database
}

func NewArticles(database *db.Database) *PvzRepo {
	return &PvzRepo{db: database}
}

func (r *PvzRepo) Add(ctx context.Context, pvz *repository.Pvz) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO pvz(pvzname,address,email) VALUES ($1,$2,$3) RETURNING id;`, pvz.PvzName, pvz.Address, pvz.Email).Scan(&id)
	return id, err
}

func (r *PvzRepo) GetByID(ctx context.Context, id int64) (*repository.Pvz, error) {
	var a repository.Pvz
	fmt.Println("pppppppp")
	err := r.db.Get(ctx, &a, `SELECT id,pvzname,address,email FROM pvz where id=$1`, id)
	if err != nil {
		fmt.Println("err")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		fmt.Println(err)
		return nil, err
	}
	return &a, nil
}

func (r *PvzRepo) Update(ctx context.Context, id int64, pvz *repository.Pvz) error {
	_, err := r.db.Exec(ctx, `UPDATE pvz SET pvzname=$1, address=$2, email=$3 WHERE id=$4`, pvz.PvzName, pvz.Address, pvz.Email, id)
	return err
}

func (r *PvzRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM pvz WHERE id=$1`, id)
	return err
}
