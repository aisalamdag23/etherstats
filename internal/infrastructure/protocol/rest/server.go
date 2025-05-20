package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aisalamdag23/etherstats/internal/infrastructure/config"
	"github.com/aisalamdag23/etherstats/internal/infrastructure/protocol/rest/middleware"
	"github.com/aisalamdag23/etherstats/internal/infrastructure/registry"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

// RunServer runs HTTP/REST server
func RunServer(ctx context.Context, cfg *config.Config, logger *zap.Logger) error {
	wait, err := time.ParseDuration(fmt.Sprintf("%ds", cfg.General.ShutdownWaitSec))
	if err != nil {
		return err
	}

	reg := registry.Init(ctx, cfg, logger)

	r := mux.NewRouter()
	r.Use(middleware.CtxWithLogger(logger))

	v1 := r.PathPrefix("/api/v1").Subrouter()

	ethServer, err := reg.CreateETHServer()
	if err != nil {
		return fmt.Errorf("failed to create eth server: %w", err)
	}

	ethServer.RegisterRoutes(v1)

	// add auth and routes here - start

	// add auth and routes here - end

	// CORS control
	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodPut},
	})
	// This inserts the middleware
	handler := cor.Handler(r)

	srv := &http.Server{
		Addr:         cfg.General.HTTPAddr,
		WriteTimeout: time.Second * time.Duration(cfg.General.WriteTimeoutSec),
		ReadTimeout:  time.Second * time.Duration(cfg.General.ReadTimeoutSec),
		IdleTimeout:  time.Second * time.Duration(cfg.General.IdleTimeoutSec),
		Handler:      handler,
	}

	// Accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		logger.Info("starting HTTP/REST server...")
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("failed to listen and serve", zap.Error(err))
		}
	}()

	// Block until the signal.
	<-c
	logger.Info("shutting down server...")
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)

	logger.Info("shutdown complete")
	os.Exit(0)

	return nil
}
