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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(h.Session.LoadAndSave)

	r.Get("/health", h.Wrap(handlers.GetHealth))

	r.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	// TODO: implement home page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/habits", http.StatusFound)
	})

	userHandler := handlers.NewUserHandler(h)
	r.Get("/users/signup", h.Wrap(userHandler.GetSignupForm))
	r.Post("/users/signup", h.Wrap(userHandler.PostSignup))
	r.Get("/users/login", h.Wrap(userHandler.GetLoginForm))
	r.Post("/users/login", h.Wrap(userHandler.PostLogin))
	r.Post("/users/logout", h.Wrap(userHandler.PostLogout))
	r.Get("/users/account", h.Wrap(userHandler.GetAccount))

	habitHandler := handlers.NewHabitHandler(h)
	r.Get("/habits", h.Wrap(habitHandler.GetHabits))
	r.Get("/habits/{id}", h.Wrap(habitHandler.GetHabit))
	r.Delete("/habits/{id}", h.Wrap(habitHandler.DeleteHabit))
	r.Get("/habits/new", h.Wrap(habitHandler.GetCreateHabitForm))
	r.Post("/habits/new", h.Wrap(habitHandler.PostHabit))
	r.Get("/habits/{id}/edit", h.Wrap(habitHandler.GetUpdateHabitForm))
	r.Post("/habits/{id}/edit", h.Wrap(habitHandler.PostUpdateHabit))

	entryHandler := handlers.NewEntryHandler(h)
	r.Post("/entries", h.Wrap(entryHandler.PostEntry))
	r.Delete("/entries/{id}", h.Wrap(entryHandler.DeleteEntry))

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
