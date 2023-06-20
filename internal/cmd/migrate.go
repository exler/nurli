package cmd

import (
	"errors"

	"github.com/exler/nurli/internal/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v2"
)

var migrateCmd = &cli.Command{
	Name:  "migrate",
	Usage: "Apply changes to the database structure",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "down",
			Usage: "Migrate all the way down",
			Value: false,
		},
		&cli.UintFlag{
			Name:  "version",
			Usage: "Migrate to a specific version",
			Value: 0,
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Force migration to a specific version",
			Value: false,
		},
	},
	Action: func(cCtx *cli.Context) error {
		db, err := openDatabase(cCtx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error opening database")
		}

		migrationOptions := database.MigrationOptions{
			Down:    cCtx.Bool("down"),
			Force:   cCtx.Bool("force"),
			Version: cCtx.Uint("version"),
		}
		if err := db.Migrate(migrationOptions); errors.Is(err, migrate.ErrNoChange) {
			logger.Info().Msg("No migrations to apply")
		} else if err != nil {
			logger.Fatal().Err(err).Msg("Error running migration")
		} else {
			if migrationOptions.Down {
				logger.Info().Msg("All migrations rolled back successfully")
			} else if migrationOptions.Version > 0 {
				logger.Info().Msgf("Migrated to version %d successfully", migrationOptions.Version)
			} else {
				logger.Info().Msg("Migrations applied successfully")
			}
		}

		return nil
	},
}
