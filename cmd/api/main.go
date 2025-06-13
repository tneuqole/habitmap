package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tneuqole/habitmap/internal/handlers"
	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/internal/session"
	"github.com/tneuqole/habitmap/internal/util"
)

const (
	readTimeout  = 10
	writeTimeout = 10
	idleTimeout  = 120
)

func main() {
	// TODO: make log level env var
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	db, err := sql.Open("sqlite3", "./habitmap.db") // TODO: probably shouldn't expose filename
	if err != nil {
		logger.Error("failed to connect to database", util.ErrorSlog(err))
		os.Exit(1)
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("error closing database connection", util.ErrorSlog(err))
			os.Exit(1)
		}
	}()

	queries := model.New(db)

	h := &handlers.BaseHandler{
		Logger:  logger,
		Queries: queries,
		Session: session.New(),
	}

	// TODO: read chi docs
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(h.Session.LoadAndSave)

	// TODO: custom timeout middleware with error page

	r.Get("/health", h.Wrap(handlers.GetHealth))

	r.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	r.Get("/", h.Wrap(h.GetHome))

	userHandler := handlers.NewUserHandler(h)
	r.Route("/users", func(r chi.Router) {
		// public routes
		r.Get("/signup", h.Wrap(userHandler.GetSignupForm))
		r.Post("/signup", h.Wrap(userHandler.PostSignup))
		r.Get("/login", h.Wrap(userHandler.GetLoginForm))
		r.Post("/login", h.Wrap(userHandler.PostLogin))

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(h.RequireAuth)

			r.Get("/account", h.Wrap(userHandler.GetAccount))
			r.Post("/logout", h.Wrap(userHandler.PostLogout))
		})
	})

	habitHandler := handlers.NewHabitHandler(h)
	r.Route("/habits", func(r chi.Router) {
		r.Use(h.RequireAuth)

		r.Get("/", h.Wrap(habitHandler.GetHabits))
		r.Get("/new", h.Wrap(habitHandler.GetCreateHabitForm))
		r.Post("/new", h.Wrap(habitHandler.PostHabit))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.Wrap(habitHandler.GetHabit))
			r.Get("/edit", h.Wrap(habitHandler.GetUpdateHabitForm))
			r.Post("/edit", h.Wrap(habitHandler.PostUpdateHabit))
			r.Delete("/", h.Wrap(habitHandler.DeleteHabit))
		})
	})

	entryHandler := handlers.NewEntryHandler(h)
	r.Route("/entries", func(r chi.Router) {
		r.Use(h.RequireAuth)

		r.Post("/new", h.Wrap(entryHandler.PostEntry))
		r.Delete("/{id}", h.Wrap(entryHandler.DeleteEntry))
	})

	// TODO: https
	logger.Info("Running on http://localhost:4000")
	srv := &http.Server{
		Addr:         ":4000",
		Handler:      r,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
		IdleTimeout:  idleTimeout * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("Error starting http server", util.ErrorSlog(err))
	}
}
