package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBURL string
}

func LoadConfig() *Config {
	// .env dosyasını yükle (WSL üzerinde projenin ana dizininde olmalı)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// PostgreSQL bağlantı cümlesini (Connection String) oluşturuyoruz
	// Format: postgres://username:password@localhost:5432/database_name?sslmode=disable
	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") +
		"/" + os.Getenv("DB_NAME") + "?sslmode=disable"

	return &Config{
		Port:  os.Getenv("PORT"),
		DBURL: dbURL,
	}
}
