package main

import (
	"orgstructure/internal/config"
	"orgstructure/internal/logger"
	"orgstructure/internal/storage"
)

func main() {
	cfg := config.New()
	log, logFile := logger.New(cfg.Logger)
	defer logFile.Close()
	log.Debug("logger initialization completed")

	db, err := storage.New(cfg)
	if err != nil {
		log.Error("database connection error", logger.Err(err))
		return
	}
	defer func() {
		dbInstance, _ := db.DB.DB()
		_ = dbInstance.Close()
	}()
	log.Debug("connection to the database is established")

	// TODO: migrations: goose

	// TODO: init router: net/http

	// TODO: run server
}
