package cmd

import (
	"github.com/exler/nurli/internal/server"
	"github.com/urfave/cli/v2"
)

var serveCmd = &cli.Command{
	Name:  "serve",
	Usage: "Run web server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "data-dir",
			EnvVars: []string{"NURLI_DATA_DIR"},
			Value:   "",
		},
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "Port to listen on",
			Value:   8000,
		},
		&cli.StringFlag{
			Name:    "username",
			EnvVars: []string{"NURLI_USERNAME"},
			Usage:   "Username for basic auth",
		},
		&cli.StringFlag{
			Name:    "password",
			EnvVars: []string{"NURLI_PASSWORD"},
			Usage:   "Password for basic auth",
		},
	},
	Action: func(cCtx *cli.Context) error {
		db, err := openDatabase(cCtx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error opening database")
		}

		port := cCtx.Int("port")

		serverConfig := server.ServerConfig{
			ServerPort:        port,
			BasicAuthUsername: cCtx.String("username"),
			BasicAuthPassword: cCtx.String("password"),
		}

		err = server.ServeApp(serverConfig, db, &logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error running server")
		}

		return nil
	},
}
