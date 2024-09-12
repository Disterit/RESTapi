package main

import (
	"RESTapi/internal/api"
	"RESTapi/internal/config"
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

	router := chi.NewRouter()
	router.Post("api/v1/wallet", api.CreateWallet)
	router.Post("api/v1/wallet/{walletID}/send", api.SendMoney)
	router.Get("/api/v1/wallet/{walletID}/history", api.HistoryWallet)
	router.Get("/api/v1/wallet/{walletID}", api.InfoWallet)

	http.ListenAndServe(cfg.Address+":"+cfg.Port, router)

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
