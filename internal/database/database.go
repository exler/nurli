package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type SQLiteDatabase struct {
	sqlx.DB
}

func OpenSQLiteDatabase(ctx context.Context, dbPath string) (sqliteDB *SQLiteDatabase, err error) {
	db, err := sqlx.ConnectContext(ctx, "sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	sqliteDB = &SQLiteDatabase{*db}
	return sqliteDB, err
}

func (db *SQLiteDatabase) Migrate() error {
	return nil
}
