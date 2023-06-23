package database

import (
	"context"

	"github.com/exler/nurli/internal"
	"github.com/exler/nurli/internal/core"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

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

func (db *SQLiteDatabase) CreateUser(ctx context.Context, user core.User) error {
	if err := db.withTx(ctx, func(tx *sqlx.Tx) error {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return err
		}

		_, err = tx.Exec("INSERT INTO users (username, password) VALUES (?, ?)",
			user.Username, hashedPassword,
		)

		return err
	}); err != nil {
		return err
	}

	return nil
}

func (db *SQLiteDatabase) ListUsers(ctx context.Context) ([]core.User, error) {
	query := `SELECT id, username, password, date_joined FROM users`
	var users []core.User
	err := db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db *SQLiteDatabase) GetUserByUsername(ctx context.Context, username string) (core.User, error) {
	query := `SELECT id, username, password, date_joined FROM users WHERE username = ?`
	var user core.User
	err := db.SelectContext(ctx, &user, query, username)
	if err != nil {
		return core.User{}, err
	}

	return user, nil
}

func (db *SQLiteDatabase) DeleteUserByUsername(ctx context.Context, username string) error {
	if err := db.withTx(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.Exec("DELETE FROM users WHERE username = ?", username)
		return err
	}); err != nil {
		return err
	}

	return nil
}

func (db *SQLiteDatabase) CreateBookmark(ctx context.Context, bookmark core.Bookmark) error {
	if err := db.withTx(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.Exec("INSERT INTO bookmarks (title, url, owner_id) VALUES (?, ?, ?)",
			bookmark.Title, bookmark.URL, bookmark.OwnerID,
		)

		return err
	}); err != nil {
		return err
	}

	return nil
}

func (db *SQLiteDatabase) ListBookmarks(ctx context.Context) ([]core.Bookmark, error) {
	query := `SELECT id, title, url, owner_id, read, favorite FROM bookmarks`
	var bookmarks []core.Bookmark
	err := db.SelectContext(ctx, &bookmarks, query)
	if err != nil {
		return nil, err
	}

	return bookmarks, nil
}

func (db *SQLiteDatabase) ToggleBookmarkRead(ctx context.Context, bookmarkID int) error {
	query := `UPDATE bookmarks SET read = (1 - read) WHERE id = ?`
	_, err := db.ExecContext(ctx, query, bookmarkID)
	if err != nil {
		return err
	}

	return nil
}

func (db *SQLiteDatabase) ToggleBookmarkFavorite(ctx context.Context, bookmarkID int) error {
	query := `UPDATE bookmarks SET favorite = (1 - favorite) WHERE id = ?`
	_, err := db.ExecContext(ctx, query, bookmarkID)
	if err != nil {
		return err
	}

	return nil
}

func (db *SQLiteDatabase) DeleteBookmark(ctx context.Context, bookmarkID int) error {
	if err := db.withTx(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.Exec("DELETE FROM bookmarks WHERE id = ?", bookmarkID)
		return err
	}); err != nil {
		return err
	}

	return nil
}
