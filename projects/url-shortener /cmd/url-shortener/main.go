package main

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"url-shortener_/internal/config"
	"url-shortener_/internal/lib/logger/sl"
	"url-shortener_/internal/storage/sqlite"
)

const (
	envLocal      = "local"
	envDev        = "development"
	envProduction = "production"
)

func main() {
	loadEnv()

	cfg := config.LoadConfig()
	log := setupLogger(cfg.Env)
	log.Info("starting..", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("storage init failed", sl.Err(err))
		os.Exit(1)
	}

	_ = storage
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProduction:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
