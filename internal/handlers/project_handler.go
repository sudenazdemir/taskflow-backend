package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sudenazdemir/taskflow-backend/internal/database"
	"github.com/sudenazdemir/taskflow-backend/internal/models"
)

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var p models.Project
	// Postman'den gelen JSON'u senin Project modeline döküyoruz
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
		return
	}

	// Veritabanı sorgusu (Senin modelindeki alanları kullanıyoruz)
	query := `INSERT INTO projects (name, description, owner_id) 
              VALUES ($1, $2, $3) RETURNING id, created_at`

	err := database.DB.QueryRow(query, p.Name, p.Description, p.OwnerID).
		Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		http.Error(w, "Proje oluşturulamadı", http.StatusInternalServerError)
		return
	}

	// Başarılı! Oluşan projeyi geri döndür
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
