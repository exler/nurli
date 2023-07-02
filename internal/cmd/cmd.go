package cmd

import (
	"os"
	"path/filepath"

	"github.com/exler/nurli/internal/database"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

var logger zerolog.Logger

var Cmd = &cli.App{
	Name:  "nurli",
	Usage: "Self-hosted and lightning-fast bookmark manager",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "data-dir",
			EnvVars: []string{"NURLI_DATA_DIR"},
			Value:   "",
		},
	},
	Commands: []*cli.Command{versionCmd, serveCmd, migrateCmd, bookmarkCmd, tagCmd, importCmd},
	Before: func(cCtx *cli.Context) error {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
		return nil
	},
}

func Execute() error {
	return Cmd.Run(os.Args)
}

func openDatabase(cCtx *cli.Context) (*gorm.DB, error) {
	var dataDir string
	var err error
	if dataDir = cCtx.String("data-dir"); dataDir == "" {
		dataDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	dbPath := filepath.Join(dataDir, "nurli.db")
	db, err := database.OpenDatabase(cCtx.Context, dbPath)
	return db, err
}
