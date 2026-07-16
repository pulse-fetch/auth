package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"auth/internal/domain"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const LocalEnv = "local"
const DevEnv = "dev"
const ProdEnv = "prod"

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	params := domain.PostgresParams{
		Host:     cfg.PostgresCfg.Host,
		Port:     cfg.PostgresCfg.Port,
		User:     cfg.PostgresCfg.User,
		DBName:   cfg.PostgresCfg.DBName,
		Password: cfg.PostgresCfg.Password,
		Sslmode:  cfg.PostgresCfg.Sslmode,
	}
	application := app.New(log, cfg.GrpcCfg.Port, params, cfg.HmacSecret)
	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	application.GRPCServer.Stop()
	log.Info("Server gracefully stopped", slog.String("signal", sign.String()))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case LocalEnv:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case DevEnv:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case ProdEnv:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log

}
