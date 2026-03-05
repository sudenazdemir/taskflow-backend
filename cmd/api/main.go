package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/sudenazdemir/taskflow-backend/internal/config"
	"github.com/sudenazdemir/taskflow-backend/internal/database"
	"github.com/sudenazdemir/taskflow-backend/internal/handlers"
	"github.com/sudenazdemir/taskflow-backend/internal/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env dosyası bulunamadı, sistem değişkenleri kullanılacak")
	}
	cfg := config.LoadConfig()

	database.Connect(cfg.DBURL)
	database.CreateTables()

	// Rotalar (Routes)
	http.HandleFunc("/user", handlers.GetUserHandler)
	http.HandleFunc("/register", middleware.LoggingMiddleware(handlers.RegisterHandler))
	http.HandleFunc("/login", middleware.LoggingMiddleware(handlers.LoginHandler))
	http.HandleFunc("/projects", middleware.LoggingMiddleware(handlers.CreateProjectHandler))
	http.HandleFunc("/projects/stats/", middleware.AuthMiddleware(handlers.GetProjectStatsHandler))
	http.HandleFunc("/tasks/", middleware.LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
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
			// Burada da JSON hata dönmek istersen az önce yazdığımız sendJSONError mantığını kullanabilirsin
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Status: OK"))
	})

	port := ":" + cfg.Port
	fmt.Println("🚀 TaskFlow Backend starting on http://localhost" + port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
