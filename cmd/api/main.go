package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sudenazdemir/taskflow-backend/internal/database"
	"github.com/sudenazdemir/taskflow-backend/internal/router"
)

func main() {
	// 1. Env yükle
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env dosyası yüklenirken hata oluştu! Dosyanın doğru yerde olduğundan emin ol.")
	}
	// 2. Veritabanı bağlantısı
	dbURL := os.Getenv("DATABASE_URL")
	database.Connect(dbURL)
	database.CreateTables()

	// 3. Router kurulumu
	routes := router.SetupRoutes()

	// 4. Server başlat
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 TaskFlow sunucusu %s portunda başladı...", port)
	log.Fatal(http.ListenAndServe(":"+port, routes))
}
