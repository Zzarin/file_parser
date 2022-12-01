package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Db struct {
	db *sqlx.DB
}

type Dbstructure struct {
	ID    int `db:"id"`
	Key   int `db:"key"`
	Value int `db:"value"`
}

func NewStorage(postgres *sqlx.DB) *Db {
	return &Db{
		db: postgres,
	}
}

func (d *Db) Init(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS local_file (
                          id serial PRIMARY KEY NOT NULL,
                          key         INT  UNIQUE NOT NULL,
                          value      INT NOT NULL);`

	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("create a database table %w", err)
	}
	return nil
}

func (d *Db) IsKeyExist(ctx context.Context, key int) (bool, error) {
	ctxDbTimeout, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	dbRespond := &Dbstructure{}
	err := d.db.GetContext(
		ctxDbTimeout,
		dbRespond,
		"SELECT * FROM local_file WHERE key = $1;",
		key,
	)
	if dbRespond == nil && err != nil {
		return false, fmt.Errorf("find data by key:#%d, %w", key, err)
	}
	if dbRespond != nil && err != nil {
		return false, nil
	}

	return true, nil
}

func (d *Db) UpdateKey(ctx context.Context, key int, value int) error {
	ctxDbTimeout, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	_, err := d.db.ExecContext(
		ctxDbTimeout,
		"UPDATE local_file SET value = $2 WHERE key = $1",
		key,
		value,
	)
	if err != nil {
		return fmt.Errorf("update key:%d, %w", key, err)
	}
	return nil
}

func (d *Db) InsertNewKey(ctx context.Context, key int, value int) error {
	ctxDbTimeout, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	_, err := d.db.ExecContext(ctxDbTimeout, `INSERT INTO local_file (key, value) VALUES ($1, $2);`, key, value)
	if err != nil {
		return fmt.Errorf("key is not recorded %w", err)
	}
	return nil
}
