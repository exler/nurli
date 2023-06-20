package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/exler/nurli/internal/database"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

type ServerConfig struct {
	DB         *database.SQLiteDatabase
	ServerPort int
	Logger     *zerolog.Logger
}

type ServerHandler struct {
	Config ServerConfig

	templates *template.Template
}

func ServeApp(config ServerConfig) error {
	sh := ServerHandler{
		Config: config,
	}
	err := sh.prepareTemplates()
	if err != nil {
		config.Logger.Fatal().Err(err).Msg("Error preparing templates")
	}

	fs := http.FileServer(http.Dir("./internal/static"))

	router := chi.NewRouter()

	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	router.Get("/", sh.IndexHandler)

	router.Route("/api", func(r chi.Router) {
		r.Get("/health", sh.HealthHandler)
	})

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
