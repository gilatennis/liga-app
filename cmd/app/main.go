package main

import (
	"errors"
	"fmt"
	"liga-app/config"
	"liga-app/infra/http"
	"liga-app/infra/logger"
	"liga-app/infra/mysql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	_, err := os.Stat(".env")
	if errors.Is(err, os.ErrNotExist) {
		log.Println("[WARNING] no .env file found")
		return nil
	}

	err = godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	return nil
}

func main() {
	err := LoadEnv()
	if err != nil {
		log.Fatalf("[cmd.app.main] LoadEnv(): failed to load env: %v", err)
	}

	// load and validate config
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config.NewConfig(): failed to load config: %v", err)
	}

	// connect to database
	db, err := mysql.NewGormDB(config)
	if err != nil {
		log.Fatalf("mysql.NewDB(): failed to connect to database: %v", err)
	}

	// setup logging
	logger, err := logger.NewLogger(config.Log.Dir)
	if err != nil {
		log.Fatalf("logger.NewLogger(): failed to setup logger: %v", err)
	}

	// setup server
	server, err := http.NewServer(config, logger, db)
	if err != nil {
		log.Fatalf("http.NewServer(): failed to setup server: %v", err)
	}

	// start server
	server.Start()
}
