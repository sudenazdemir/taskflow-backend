package router

import (
	"net/http"

	"github.com/sudenazdemir/taskflow-backend/internal/handlers"
	"github.com/sudenazdemir/taskflow-backend/internal/middleware"
)

// router.go içine şu fonksiyonu ekle
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Şimdilik herkes (Test için)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// --- Public Routes (Giriş Gerektirmeyenler) ---
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)

	// --- Protected Routes (JWT Gerektirenler) ---
	// Projeler
	mux.HandleFunc("/projects/stats/", middleware.AuthMiddleware(handlers.GetProjectStatsHandler))

	// Görevler ve Yorumlar
	mux.HandleFunc("/comments/add", middleware.AuthMiddleware(handlers.AddCommentHandler))
	mux.HandleFunc("/tasks/details/", middleware.AuthMiddleware(handlers.GetTaskWithCommentsHandler))

	// Projeye ait TÜM görevleri listelemek için:
	mux.HandleFunc("/projects/tasks/", middleware.AuthMiddleware(handlers.GetTasksHandler))
	mux.HandleFunc("/tasks/add", middleware.AuthMiddleware(handlers.CreateTaskHandler))
	mux.HandleFunc("/tasks/update/", middleware.AuthMiddleware(handlers.UpdateTaskHandler))
	mux.HandleFunc("/tasks/delete/", middleware.AuthMiddleware(handlers.DeleteTaskHandler))

	mux.HandleFunc("/tasks/upload", middleware.AuthMiddleware(handlers.UploadFileHandler))
	mux.HandleFunc("/tasks/attachments", middleware.AuthMiddleware(handlers.GetAttachmentsHandler))

	return enableCORS(mux) // CORS ile sarmalayıp gönder
}
