package database

import "github.com/jmoiron/sqlx"

// Interface for accessing and manipulating data in database.
type Database interface {
	Migrate() error
}

type DatabaseBase struct {
	sqlx.DB
}
