package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Dün Network dersinde öğrendiğimiz Port kavramı!
	port := ":8080"

	// Basit bir Health Check (Sağlık Kontrolü) endpoint'i
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Cevabın JSON formatında olacağını belirtiyoruz (Header)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// O çok sevdiğin "milisaniyelik" hız testi için zamanı da ekleyelim
		response := map[string]string{
			"status":  "success",
			"message": "Hello TaskFlow! The backend is running perfectly.",
			"time":    time.Now().Format(time.RFC3339),
		}

		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("🚀 TaskFlow Backend starting on http://localhost" + port)

	// Sunucuyu ayağa kaldırıyoruz
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
