package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sudenazdemir/taskflow-backend/internal/database"
	"github.com/sudenazdemir/taskflow-backend/internal/models"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Veritabanından veriyi çek (Şimdilik id=1 sabit, sonra dinamik yapacağız)
	var user models.User
	err := database.DB.QueryRow("SELECT id, username, email FROM users WHERE id = 1").Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}

	// 2. Veriyi JSON formatına çevir ve gönder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
