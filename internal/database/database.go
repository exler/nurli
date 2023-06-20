package database

import (
	"context"

	"github.com/exler/nurli/internal"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
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

func (db *SQLiteDatabase) Migrate(options MigrationOptions) error {
	sourceDriver, err := iofs.New(internal.MigrationsFS, "database/migrations")
	if err != nil {
		return err
	}

	dbDriver, err := sqlite.WithInstance(db.DB.DB, &sqlite.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite", dbDriver)
	if err != nil {
		return err
	}

	if options.Down {
		err = migration.Down()
	} else if options.Version > 0 && options.Force {
		err = migration.Force(int(options.Version))
	} else if options.Version > 0 {
		err = migration.Migrate(options.Version)
	} else {
		err = migration.Up()
	}

	return err
}
