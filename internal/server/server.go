package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/exler/nurli/internal"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ServerConfig struct {
	ServerPort        int
	BasicAuthUsername string
	BasicAuthPassword string
}

type ServerHandler struct {
	DB     *gorm.DB
	Logger *zerolog.Logger

	templates *template.Template
}

func ServeApp(config ServerConfig, db *gorm.DB, logger *zerolog.Logger) error {
	sh := ServerHandler{
		DB:     db,
		Logger: logger,
	}
	err := sh.prepareTemplates()
	if err != nil {
		logger.Fatal().Err(err).Msg("Error preparing templates")
	}

	fs := http.FileServer(http.FS(internal.StaticFS))

	router := chi.NewRouter()
	if config.BasicAuthUsername != "" && config.BasicAuthPassword != "" {
		router.Use(BasicAuthMiddleware("Nurli", config.BasicAuthUsername, config.BasicAuthPassword))
	}
	router.Handle("/static/*", fs)

	// UI routes
	router.Route("/", func(r chi.Router) {
		r.Get("/", sh.IndexHandler)

		r.Get("/add", sh.AddBookmarkHandler)
		r.Post("/add", sh.AddBookmarkHandler)

		r.Get("/edit/{id}", sh.EditBookmarkHandler)
		r.Post("/edit/{id}", sh.EditBookmarkHandler)

		r.Get("/delete/{id}", sh.DeleteBookmarkHandler)
		r.Post("/delete/{id}", sh.DeleteBookmarkHandler)
	})

	// API routes
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

	logger.Info().Msgf("Listening on port %d", config.ServerPort)

	return srv.ListenAndServe()
}
