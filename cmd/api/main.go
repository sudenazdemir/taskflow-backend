package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sudenazdemir/taskflow-backend/internal/config"
	"github.com/sudenazdemir/taskflow-backend/internal/database"
)

func main() {
	// 1. Ayarları (.env) yükle
	cfg := config.LoadConfig()

	// 2. Veritabanına bağlan
	database.Connect(cfg.DBURL)

	// 3. Sunucuyu başlat
	port := ":" + cfg.Port
	fmt.Println("🚀 TaskFlow Backend starting on http://localhost" + port)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Status: OK"))
	})

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
