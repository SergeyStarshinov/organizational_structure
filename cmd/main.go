package main

import (
	"orgstructure/internal/config"
	"orgstructure/internal/logger"
)

func main() {
	cfg := config.New()
	log := logger.New(cfg.Logger)
	log.Debug("logger initialization complete")

	// TODO: connect DB: postgreSQL

	// TODO: migrations: goose

	// TODO: init router: net/http

	// TODO: run server
}
