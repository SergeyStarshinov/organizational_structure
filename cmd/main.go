package main

import (
	"net/http"
	"orgstructure/internal/config"
	"orgstructure/internal/logger"
	"orgstructure/internal/storage"
	"orgstructure/internal/web"
)

func main() {

	cfg := config.New()
	log, logFile := logger.New(cfg.Logger)
	defer logFile.Close()
	log.Debug("logger initialization completed")

	repo, err := storage.New(cfg)
	if err != nil {
		log.Error("database connection error", logger.Err(err))
		return
	}
	defer func() {
		dbInstance, _ := repo.DB.DB()
		_ = dbInstance.Close()
	}()
	log.Debug("connection to the database is established")

	handler := web.NewBaseHandler(repo, log)
	mux := web.CreateMux(handler)

	server := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: mux,
	}
	server.ListenAndServe()
	log.Info("server started")

}
