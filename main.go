package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/NessibeliY/quotes-api/config"
	"github.com/NessibeliY/quotes-api/internal/handlers"
	"github.com/NessibeliY/quotes-api/internal/service"
	"github.com/NessibeliY/quotes-api/internal/store"
	"github.com/NessibeliY/quotes-api/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	l, err := logger.NewLogger(cfg.LogFile)
	if err != nil {
		slog.Error("load logger", "error", err)
		os.Exit(1)
	}
	slog.SetDefault(l)
	l.Info("Logger loaded", "file", cfg.LogFile)

	store := store.NewStore()
	service := service.NewService(store)
	handler := handlers.NewHandler(service)

	r1 := handlers.NewRateLimiter(100, time.Minute)

	r := mux.NewRouter()
	r.Use(r1.LimitMiddleware)

	r.HandleFunc("/quotes", handler.AddQuote).Methods("POST")
	r.HandleFunc("/quotes", handler.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", handler.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id:[0-9]+}", handler.DeleteQuote).Methods("DELETE")
	r.HandleFunc("/health", handler.Health).Methods("GET")

	ctx, cancel := context.WithCancel(context.Background())

	errChan := make(chan error, 1)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		slog.Info("shutdown signal received")
		cancel()
	}()

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	go func() {
		slog.Info("server started", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("server shutdown error", "error", err)
		} else {
			slog.Info("server gracefully stopped")
		}
	case err := <-errChan:
		slog.Error("server error", "error", err)
	}
}
