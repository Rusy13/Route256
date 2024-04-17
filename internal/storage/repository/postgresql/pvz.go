package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"log"

	"Homework/internal/storage/db"
	"Homework/internal/storage/repository"
)

type PvzRepo struct {
	db db.PGX
}

func NewPvzRepo(database db.PGX) *PvzRepo {
	return &PvzRepo{db: database}
}

func (r *PvzRepo) Add(ctx context.Context, pvz *repository.Pvz) (int64, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	var id int64
	err = tx.QueryRow(ctx, `INSERT INTO pvz(pvzname,address,email) VALUES ($1,$2,$3) RETURNING id;`, pvz.PvzName, pvz.Address, pvz.Email).Scan(&id)

	if err != nil {
		log.Println("Error in Add!!", err)
	}

	log.Println("Committing transaction...")
	if err := tx.Commit(ctx); err != nil {
		return id, err
	}

	log.Println("Transaction committed successfully!")
	return id, err
}

func (r *PvzRepo) GetByID(ctx context.Context, id int64) (*repository.Pvz, error) {
	var a repository.Pvz
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return &a, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = tx.QueryRow(ctx, `SELECT id, pvzname, address, email FROM pvz WHERE id = $1`, id).Scan(&a.ID, &a.PvzName, &a.Address, &a.Email)
	if err != nil {
		log.Println("Error in GetByID!!", err)
	}

	log.Println("Committing transaction...")
	if err := tx.Commit(ctx); err != nil {
		return &a, err
	}

	log.Println("Transaction committed successfully!")
	return &a, err

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *PvzRepo) Update(ctx context.Context, id int64, pvz *repository.Pvz) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `UPDATE pvz SET pvzname=$1, address=$2, email=$3 WHERE id=$4`, pvz.PvzName, pvz.Address, pvz.Email, id)
	if err != nil {
		log.Println("Error in Update!!", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil

}

func (r *PvzRepo) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `DELETE FROM pvz WHERE id=$1`, id)
	if err != nil {
		log.Println("Error in Delete!!", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
