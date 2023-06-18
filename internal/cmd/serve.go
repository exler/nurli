package cmd

import (
	"github.com/exler/nurli/internal/server"
	"github.com/urfave/cli/v2"
)

var serveCmd = &cli.Command{
	Name:  "serve",
	Usage: "Run web server",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "Port to listen on",
			Value:   8000,
		},
	},
	Action: func(cCtx *cli.Context) error {
		db, err := openDatabase(cCtx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error opening database")
		}

		port := cCtx.Int("port")

		serverConfig := server.Config{
			DB:         db,
			ServerPort: port,
			Logger:     &logger,
		}

		err = server.ServeApp(serverConfig)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error running server")
		}

		return nil
	},
}
