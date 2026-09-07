package main

import (
	"fmt"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/config"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/config/database"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/config/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize new logger
	log := logger.New()
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

	log.Info().Msg("Server is running on " + cfg.Server.Host)
}
