package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/exler/nurli/internal/core"
	"github.com/exler/nurli/internal/database"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type contextKey int

const (
	LoggerKey contextKey = iota
)

var Cmd = &cli.App{
	Name:  "nurli",
	Usage: "Self-hosted and lightning-fast bookmark manager",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "data-dir",
			EnvVars: []string{"NURLI_DATA_DIR"},
			Value:   "",
		},
		&cli.BoolFlag{
			Name:    "debug",
			EnvVars: []string{"NURLI_DEBUG"},
			Value:   false,
		},
	},
	Commands: []*cli.Command{versionCmd, serveCmd, migrateCmd, bookmarkCmd, tagCmd, importCmd},
	Before: func(cCtx *cli.Context) error {
		newLogger := core.NewZerologGORMLogger(cCtx.Bool("debug"), logger.Config{
			SlowThreshold: time.Second,
		})
		newLogger.LogMode(logger.Info)
		ctx := context.WithValue(cCtx.Context, LoggerKey, newLogger)
		cCtx.Context = ctx
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
	logger := cCtx.Context.Value(LoggerKey).(*core.ZerologGORMLogger)
	db, err := database.OpenDatabase(logger, dbPath)
	return db, err
}
