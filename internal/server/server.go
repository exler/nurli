package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/exler/nurli/internal/database"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

type Config struct {
	DB         database.Database
	ServerPort int
	Logger     *zerolog.Logger
}

func ServeApp(config Config) error {
	router := chi.NewRouter()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.ServerPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	config.Logger.Info().Msgf("Listening on port %d", config.ServerPort)

	return srv.ListenAndServe()
}
