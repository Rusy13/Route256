package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"log"
	"strconv"
	"time"

	"Homework/internal/storage/db"
	"Homework/internal/storage/repository"
	InMemoryCache "Homework/internal/storage/repository/in_memory_cache"
)

const cacheTTL = time.Second * 15

type PvzRepo struct {
	db          db.PGX
	queryCache  *InMemoryCache.InMemoryCache // Кеш для запросов к базе данных
	resultCache *InMemoryCache.InMemoryCache // Кеш для результатов запросов
}

func NewPvzRepo(database db.PGX) *PvzRepo {
	return &PvzRepo{
		db:          database,
		queryCache:  InMemoryCache.NewInMemoryCache(),
		resultCache: InMemoryCache.NewInMemoryCache(),
	}
}

func (r *PvzRepo) startCacheCleanup() {
	go func() {
		ticker := time.NewTicker(time.Second * 5) // Периодический интервал сканирования кеша
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.resultCache.CleanupExpiredEntries() // Вызываем функцию очистки истекших записей из кеша
			}
		}
	}()
}

func (r *PvzRepo) Add(ctx context.Context, pvz *repository.Pvz) (int64, error) {
	r.resultCache.Clear()

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
	cacheKey := "GetByID:" + strconv.FormatInt(id, 10)
	if cachedData, found := r.resultCache.Get(cacheKey); found {
		log.Println("Data found in cache for GetByID operation")
		if cachedData != nil {
			return cachedData.(*repository.Pvz), nil
		}
		return nil, nil
	}

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
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No rows found!")
			r.resultCache.Set(cacheKey, nil, cacheTTL)
			return nil, nil
		} else {
			log.Printf("Error occurred while querying data: %v", err)
			return nil, err
		}
	}

	if err != nil {
		log.Printf("Error occurred: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No rows found, nothing to cache")
		}
		return nil, err
	} else {
		r.resultCache.Set(cacheKey, &a, cacheTTL)
		log.Println("Data cached successfully")
	}

	log.Println("Committing transaction...")
	if err := tx.Commit(ctx); err != nil {
		return &a, err
	}

	log.Println("Transaction committed successfully!")
	return &a, err
}

func (r *PvzRepo) Update(ctx context.Context, id int64, pvz *repository.Pvz) error {
	r.resultCache.Clear()

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
	r.resultCache.Clear()

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
