package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sudenazdemir/taskflow-backend/internal/database"
	"github.com/sudenazdemir/taskflow-backend/internal/models"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	// 1. JWT'den gelen kullanıcıyı al
	userID := r.Context().Value("user_id").(int)

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	// 2. Veritabanına kaydet
	query := `
        INSERT INTO comments (task_id, user_id, content) 
        VALUES ($1, $2, $3) 
        RETURNING id, created_at`

	err := database.DB.QueryRow(query, comment.TaskID, userID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		http.Error(w, "Yorum eklenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	comment.UserID = userID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}
