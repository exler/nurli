package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ServerConfig struct {
	ServerPort int
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

	fs := http.FileServer(http.Dir("./internal/static"))

	router := chi.NewRouter()

	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Auth routes
	router.Group(func(r chi.Router) {
		r.Post("/login", sh.LoginHandler)
		r.Get("/login", sh.LoginHandler)
		r.Get("/logout", sh.LogoutHandler)
	})

	// UI routes
	router.Route("/", func(r chi.Router) {
		r.Use(sh.AuthMiddleware)
		r.Get("/", sh.IndexHandler)

		r.Get("/add", sh.AddBookmarkHandler)
		r.Post("/add", sh.AddBookmarkHandler)
	})

	// API routes
	router.Route("/api", func(r chi.Router) {
		r.Use(sh.AuthMiddleware)
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
