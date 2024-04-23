package app

import (
	"log/slog"

	"github.com/elusiv0/medods_test/pkg/httpserver"
)

type App struct {
	server *httpserver.HttpServer
	logger *slog.Logger
}

func New(
	serv *httpserver.HttpServer,
	log *slog.Logger,
) *App {
	app := &App{
		server: serv,
		logger: log,
	}

	return app
}

func (a *App) Run() error {
	if err := a.server.Start(); err != nil {
		a.logger.Error("error with up httpserver: " + err.Error())
		return err
	}

	return nil
}
