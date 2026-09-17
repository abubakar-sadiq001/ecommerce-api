package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/config"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/database"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/server"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	// initialize a new logger
	log := zerolog.New(os.Stdout)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
		fmt.Println("logger Err")
	}

	// connect to database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
		fmt.Println("database Err")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
		fmt.Println("mainDB Err")
	}

	defer mainDB.Close()
	gin.SetMode(cfg.Server.GinMode)

	authService := services.NewAuthService(db, cfg)
	productService := services.NewProductServices(db)
	userService := services.NewUserService(db)

	srv := server.New(cfg, db, &log, authService, productService, userService)

	// initialize parameters for running an HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      srv.SetupRoutes(), // register the routes
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// server listens and run in a goroutine
	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("starting server")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // waiting to receive shutting down signal to unblock

	log.Info().Msg("Shuting down server ")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
		return
	}

	log.Info().Err(err).Msg("shutting down database")
}
