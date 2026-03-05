package handlers

import (
	"encoding/json"
	"log"
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

func GetProjectStatsHandler(w http.ResponseWriter, r *http.Request) {
	// URL'den ID alıyoruz (getIDFromURL fonksiyonunu kullanıyoruz)
	projectID, err := getIDFromURL(r)
	if err != nil || projectID == 0 {
		sendJSONError(w, "Geçersiz Proje ID", http.StatusBadRequest)
		return
	}

	var total, completed, pending int

	// Tek sorgu, dev hizmet: Toplam, Tamamlanan ve Bekleyen sayılarını alıyoruz
	query := `
        SELECT 
            COUNT(*) as total,
            COUNT(*) FILTER (WHERE status = 'completed') as completed,
            COUNT(*) FILTER (WHERE status = 'pending') as pending
        FROM tasks 
        WHERE project_id = $1`

	err = database.DB.QueryRow(query, projectID).Scan(&total, &completed, &pending)
	if err != nil {
		log.Printf("İstatistik sorgu hatası: %v", err)
		sendJSONError(w, "İstatistikler hesaplanamadı", http.StatusInternalServerError)
		return
	}

	// Başarı oranını hesapla (0'a bölünme hatasını engellemek için)
	completionRate := 0.0
	if total > 0 {
		completionRate = (float64(completed) / float64(total)) * 100
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"project_id":      projectID,
		"total_tasks":     total,
		"completed_tasks": completed,
		"pending_tasks":   pending,
		"completion_rate": completionRate,
	})
}
