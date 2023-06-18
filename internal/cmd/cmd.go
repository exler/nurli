package cmd

import (
	"os"
	"path/filepath"

	"github.com/exler/nurli/internal/database"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var logger zerolog.Logger

var Cmd = &cli.App{
	Name:  "nurli",
	Usage: "Efficient bookmark manager",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "data-dir",
			EnvVars: []string{"DATA_DIR"},
			Value:   "",
		},
	},
	Commands: []*cli.Command{versionCmd, serveCmd, migrateCmd},
	Before: func(cCtx *cli.Context) error {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
		return nil
	},
}

func Execute() error {
	return Cmd.Run(os.Args)
}

func openDatabase(cCtx *cli.Context) (database.Database, error) {
	switch dbms := cCtx.String("NURLI_DB_TYPE"); dbms {
	default:
		return openSQLiteDatabase(cCtx)
	}
}

func openSQLiteDatabase(cCtx *cli.Context) (*database.SQLiteDatabase, error) {
	var dataDir string
	var err error
	if dataDir = cCtx.String("DATA_DIR"); dataDir == "" {
		dataDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	dbPath := filepath.Join(dataDir, "nurli.db")
	db, err := database.OpenSQLiteDatabase(cCtx.Context, dbPath)
	return db, err
}
