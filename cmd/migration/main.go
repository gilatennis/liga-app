package main

import (
	"liga-app/config"
	mysqlDB "liga-app/infra/mysql"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("[server.SetupDependencies]: error getting config: %v", err)
	}

	db, err := mysqlDB.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("[server.SetupDependencies]: error getting db: %v", err)
	}

	err = mysqlDB.GormMigrate(db)
	if err != nil {
		log.Fatalf("[server.SetupDependencies]: error migrating db: %v", err)
	}

	log.Println("DB migrated successfully")
}
