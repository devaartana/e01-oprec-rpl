package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/devaartana/e01-oprec-rpl/internal/auth"
	"github.com/devaartana/e01-oprec-rpl/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)


type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	authenticator auth.Authenticator
}

type config struct {
	addr    string
	db      dbConfig
	auth    authConfig
}

type dbConfig struct {
	addr              string
	maxOpenConnection int
	maxIdleConnection int
	maxIdleTime       string
}

type authConfig struct {
	user   string
	secret string
	exp    time.Duration
	iss    string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middleware global
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.HandleFunc("/{slug}", app.SlugHandler)
	r.Route("/api", func(r chi.Router) {
		r.Post("/register", app.RegisterUserHandler)
		r.Post("/login", app.LoginUserHandler)
		r.Get("/user", app.UserHandler)

		r.Route("/links", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Post("/", app.CreateLinkHandler)
			r.Get("/", app.GetAllLinksHandler)
			r.Put("/", app.UpdateLinkHandler)
			r.Delete("/{slug}", app.DeleteLinkHandler)
			r.Get("/refresh/{slug}", app.RefreshExpiredDateHandler)
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("server has started", "addr", app.config.addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr)

	return nil
}