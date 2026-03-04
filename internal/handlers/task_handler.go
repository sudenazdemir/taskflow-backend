package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sudenazdemir/taskflow-backend/internal/database"
	"github.com/sudenazdemir/taskflow-backend/internal/models"
)

// Yardımcı fonksiyon: Hataları JSON formatında göndermek için
func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func getIDFromURL(r *http.Request) (int, error) {
	// URL Path: /tasks/3
	// TrimPrefix ile başındaki /tasks/ kısmını atıyoruz, geriye sadece "3" kalıyor
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")

	// Eğer sonda hala "/" varsa onu da temizle (Örn: /tasks/3/)
	path = strings.Trim(path, "/")

	if path == "" {
		return 0, nil
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		log.Printf("ID Dönüştürme Hatası: Path '%s' sayıya çevrilemedi", path)
		return 0, err
	}
	return id, nil
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONError(w, "Sadece POST metodu destekleniyor", http.StatusMethodNotAllowed)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendJSONError(w, "Geçersiz JSON verisi", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO tasks (title, description, status, project_id, assigned_to) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	err := database.DB.QueryRow(query, task.Title, task.Description, "pending", task.ProjectID, task.AssignedTo).
		Scan(&task.ID, &task.CreatedAt)
	if err != nil {
		log.Printf("Veritabanı kayıt hatası: %v", err)
		sendJSONError(w, "Görev veritabanına kaydedilemedi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Görev başarıyla oluşturuldu",
		"task":    task,
	})
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONError(w, "Sadece GET metodu destekleniyor", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query("SELECT id, title, description, status, project_id, assigned_to, created_at FROM tasks")
	if err != nil {
		sendJSONError(w, "Görevler listelenirken bir hata oluştu", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task = []models.Task{} // Boşsa null değil [] dönsün diye

	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.ProjectID, &t.AssignedTo, &t.CreatedAt); err != nil {
			log.Printf("Satır okuma hatası: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(tasks),
		"tasks": tasks,
	})
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 1. URL'den ID'yi al
	id, err := getIDFromURL(r)
	if err != nil || id == 0 {
		sendJSONError(w, "Geçersiz ID formatı URL'de olmalı: /tasks/1", http.StatusBadRequest)
		return
	}
	// 2. Body'den yeni verileri al (ID göndermesine gerek yok artık)
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendJSONError(w, "Geçersiz veri formatı", http.StatusBadRequest)
		return
	}

	// URL'den gelen ID'yi modele set et (Body'de gelse bile URL'deki esastır)
	task.ID = id

	// 3. Veritabanında güncelleme yap
	query := `UPDATE tasks SET title=$1, description=$2, status=$3 
              WHERE id=$4 RETURNING project_id, assigned_to, created_at`

	err = database.DB.QueryRow(query, task.Title, task.Description, task.Status, task.ID).
		Scan(&task.ProjectID, &task.AssignedTo, &task.CreatedAt)
	if err != nil {
		log.Printf("Güncelleme hatası: %v", err)
		sendJSONError(w, "Güncellenecek görev bulunamadı", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Görev başarıyla güncellendi",
		"task":    task,
	})
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 1. URL'den ID'yi al
	id, err := getIDFromURL(r)
	if err != nil || id == 0 {
		sendJSONError(w, "Geçersiz ID formatı URL'de olmalı: /tasks/1", http.StatusBadRequest)
		return
	}

	// 2. Veritabanından sil
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		sendJSONError(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		sendJSONError(w, "Silinecek görev bulunamadı", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Görev başarıyla silindi",
		"id":      id,
	})
}
