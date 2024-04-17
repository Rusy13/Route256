package posgresql

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"Homework/internal/config"
	"Homework/internal/storage/db"
)

type TDB struct {
	DB db.DBops
	sync.Mutex
}

func NewFromEnv(c config.StorageConfig) *TDB {
	db, err := db.NewDb(context.Background(), c)
	if err != nil {
		panic(err)
	}
	return &TDB{DB: db}
}

func (d *TDB) SetUp(t *testing.T, args ...interface{}) {
	t.Helper()
	d.Lock()
	d.Truncate(context.Background())
}

func (d *TDB) TearDown() {
	defer d.Unlock()
	d.Truncate(context.Background())
}

func (d *TDB) Truncate(ctx context.Context) {
	var tables []string
	err := d.DB.Select(ctx, &tables, "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'goose_db_version'")
	if err != nil {
		panic(err)
	}
	if len(tables) == 0 {
		panic("run migration plz")
	}

	_, err = d.DB.Exec(ctx, "ALTER SEQUENCE pvz_id_seq RESTART WITH 1")
	if err != nil {
		panic(err)
	}

	q := fmt.Sprintf("Truncate table %s", strings.Join(tables, ","))
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}

}
