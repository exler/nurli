package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type MigrationOptions struct {
	// Migrate all the way down instead of up.
	Down bool

	// Sets a migration version. It does not check any currently active version in database. It resets the dirty state to false.
	Force bool

	// The version to migrate to. If 0, migrate to the latest version.
	Version uint
}

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

func (db *SQLiteDatabase) withTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Commit(); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg(fmt.Sprintf("Error during commit: %s", err.Error()))
		}
	}()

	err = fn(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg(fmt.Sprintf("Error during rollback: %s", err.Error()))
		}
		return err
	}

	return err
}
