package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sudenazdemir/taskflow-backend/internal/config"
	"github.com/sudenazdemir/taskflow-backend/internal/database"
	"github.com/sudenazdemir/taskflow-backend/internal/handlers"
)

func main() {
	cfg := config.LoadConfig()

	database.Connect(cfg.DBURL)
	database.CreateTables()

	// Rotalar (Routes)
	http.HandleFunc("/user", handlers.GetUserHandler)

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateTaskHandler(w, r)
		case http.MethodGet:
			handlers.GetTasksHandler(w, r)
		case http.MethodPut:
			handlers.UpdateTaskHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteTaskHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Status: OK"))
	})

	port := ":" + cfg.Port
	fmt.Println("🚀 TaskFlow Backend starting on http://localhost" + port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
