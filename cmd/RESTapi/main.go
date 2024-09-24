package main

import (
	"RESTapi/internal/api"
	"RESTapi/internal/config"
	"RESTapi/internal/storage"
	"RESTapi/internal/storage/postgres"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting api wallet", slog.String("key", cfg.Env))
	log.Debug("debug message are enable")

	db := storage.Connection()
	defer db.Close()

	storageDB := postgres.NewStorage(db)

	storageDB.CreateTable()

	router := chi.NewRouter()
	router.Post("/api/v1/wallet", api.CreateWalletHandler(storageDB))
	router.Post("/api/v1/wallet/{walletID}/send", api.SendMoneyHandler(storageDB))
	router.Get("/api/v1/wallet/{walletID}/history", api.HistoryWalletHandler(storageDB))
	router.Get("/api/v1/wallet/{walletID}", api.InfoWalletHandler(storageDB))

	err := http.ListenAndServe(cfg.Address, router)
	if err != nil {
		log.Error("error starting http server ", err)
	}

}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
