package app

import (
	"database/sql"
	"errors"
	"fmt"
	"lamoda-tech-assigment/internal/config"
	entryHttp "lamoda-tech-assigment/internal/entrypoints/http"
	"lamoda-tech-assigment/internal/services/uow"
	"lamoda-tech-assigment/internal/services/usecases"
)

type Server interface {
	ListenAndServe() error
}

type App struct {
	httpServer Server
	cfg        config.App
	db         *sql.DB
}

// MustNew create app instance and finalize func
func MustNew(cfg config.App) (*App, func()) {

	db, err := newPostgresqlConnect(cfg.Database)
	if err != nil {
		panic("app: can't connect to database")
	}
	runMigrations(cfg.Database)

	unitOfWork := uow.New(db)
	uc := usecases.New(unitOfWork)

	srv := entryHttp.NewServer(cfg.HTTP, uc)

	app := &App{
		srv,
		cfg,
		db,
	}

	return app, app.Finalize
}

func (a *App) Run() error {
	defer a.db.Close()

	if err := a.httpServer.ListenAndServe(); err != nil {
		return errors.New(fmt.Sprintf("app: %v", err))
	}

	return nil
}

func (a *App) Finalize() {
	_ = a.db.Close()
}
