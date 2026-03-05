package handlers

import (
	"database/sql"
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
	// URL'yi parçalara ayır (Örn: /projects/stats/2 -> ["", "projects", "stats", "2"])
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) < 1 {
		return 0, nil
	}

	// URL'nin en sonundaki parçayı ID olarak al
	lastPart := parts[len(parts)-1]

	id, err := strconv.Atoi(lastPart)
	if err != nil {
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

	query := `INSERT INTO tasks (title, description, status, project_id, assigned_to, priority, due_date) 
          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`

	// Scan kısmında da t.Status'u modelden alalım (Boş gelirse DB default'u kullanır)
	status := task.Status
	if status == "" {
		status = "pending"
	}

	err := database.DB.QueryRow(query, task.Title, task.Description, status, task.ProjectID, task.AssignedTo, task.Priority, task.DueDate).
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
	// URL'den project_id parametresini oku (Örn: /tasks?project_id=2)
	projectIDStr := r.URL.Query().Get("project_id")

	var rows *sql.Rows
	var err error

	if projectIDStr != "" {
		// Filtreli Sorgu
		projectID, _ := strconv.Atoi(projectIDStr)
		query := `SELECT id, project_id, assigned_to, title, description, status, priority, due_date, created_at 
                  FROM tasks WHERE project_id = $1`
		rows, err = database.DB.Query(query, projectID)
	} else {
		// Tümünü Getir (Eski halimiz)
		query := `SELECT id, project_id, assigned_to, title, description, status, priority, due_date, created_at FROM tasks`
		rows, err = database.DB.Query(query)
	}

	if err != nil {
		sendJSONError(w, "Sorgu hatası", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task = []models.Task{}

	for rows.Next() {
		var t models.Task
		// 2. BURAYI GÜNCELLE: Scan sırası SELECT sırasıyla birebir aynı olmalı!
		err := rows.Scan(
			&t.ID,
			&t.ProjectID,
			&t.AssignedTo,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.DueDate,
			&t.CreatedAt,
		)
		if err != nil {
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
	query := `UPDATE tasks SET title=$1, description=$2, status=$3, priority=$4, due_date=$5 
          WHERE id=$6 RETURNING project_id, assigned_to, created_at`

	err = database.DB.QueryRow(query, task.Title, task.Description, task.Status, task.Priority, task.DueDate, task.ID).
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
