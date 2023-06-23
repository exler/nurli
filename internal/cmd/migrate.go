package cmd

import (
	"fmt"

	"github.com/exler/nurli/internal/database"
	"github.com/urfave/cli/v2"
)

var migrateCmd = &cli.Command{
	Name:  "migrate",
	Usage: "Apply changes to the database structure",
	Action: func(cCtx *cli.Context) error {
		db, err := openDatabase(cCtx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error opening database")
		}

		err = db.AutoMigrate(&database.User{}, &database.Bookmark{}, &database.Tag{})
		if err != nil {
			logger.Fatal().Err(err).Msg("Error migrating database")
		} else {
			fmt.Println("Database migrated successfully")
		}

		return nil
	},
}
