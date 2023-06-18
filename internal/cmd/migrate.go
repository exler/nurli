package cmd

import (
	"github.com/urfave/cli/v2"
)

var migrateCmd = &cli.Command{
	Name:  "migrate",
	Usage: "Apply changes to the database structure",
	Flags: []cli.Flag{},
	Action: func(cCtx *cli.Context) error {
		db, err := openDatabase(cCtx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error opening database")
		}

		if err := db.Migrate(); err != nil {
			logger.Fatal().Err(err).Msg("Error running migration")
		}

		return nil
	},
}
