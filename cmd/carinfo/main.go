package main

import (
	"carinfo/internal/config"
	"carinfo/internal/repository/postgresql"
	"carinfo/internal/service"
	"carinfo/internal/transport/web"
	"carinfo/internal/transport/web/handler"
	"carinfo/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
)

func main() {
	logger.Log.Info("getting configs")
	cfg := config.All

	logger.Log.Info("creating db")
	db, err := postgresql.New(cfg.PostgreSQLDSN)
	if err != nil {
		logger.Log.Error("db creating error", err.Error(), "!")
	}
	m, err := migrate.New(
		cfg.MigrationsPath,
		cfg.PostgreSQLDSN)
	if err != nil {
		logger.Log.Error("migrations error:", err.Error(), "!")
	}
	if err := m.Up(); err != nil {
		logger.Log.Info("migrations error:", err.Error(), "!")
	}

	logger.Log.Info("car service creating")
	carService := service.NewCarService(db)

	logger.Log.Info("handlers creating")
	newHandler := handler.NewHandler(carService)

	logger.Log.Info("handlers initializing")
	handlers := newHandler.Init()
	logger.Log.Info("server creating")
	s := web.NewServer(cfg.HTTPHost)
	logger.Log.Info("starting server")
	err = s.Run(handlers)
	if err != nil {
		logger.Log.Error("cannot run the server", err.Error(), "!")
		os.Exit(2)
	}

}
