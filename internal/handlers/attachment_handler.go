package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/sudenazdemir/taskflow-backend/internal/database"
	"github.com/sudenazdemir/taskflow-backend/internal/models"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Dosyayı yakala (Maksimum 5MB limit)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Error(w, "Dosya çok büyük veya geçersiz", http.StatusBadRequest)
		return
	}

	// Formdan "attachment" key'i ile dosyayı alıyoruz
	file, handler, err := r.FormFile("attachment")
	if err != nil {
		http.Error(w, "Dosya bulunamadı", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Task ID'yi formdan alıyoruz (Örn: "12")
	taskIDStr := r.FormValue("task_id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Geçersiz Görev (Task) ID", http.StatusBadRequest)
		return
	}

	// 2. Güvenlik Kontrolü: Uzantı kontrolü yapalım
	allowedExtensions := map[string]bool{".pdf": true, ".txt": true, ".docx": true, ".go": true, ".png": true, ".jpg": true}
	ext := filepath.Ext(handler.Filename)
	if !allowedExtensions[ext] {
		http.Error(w, "Bu dosya tipine izin verilmiyor", http.StatusBadRequest)
		return
	}

	// 3. Klasör yoksa oluştur
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Klasör oluşturulamadı", http.StatusInternalServerError)
		return
	}

	// 4. Senin Kararın: UUID ile benzersiz isim veriyoruz
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, newFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Dosya oluşturulamadı", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 5. İçeriği kopyala (Diske yazma)
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Dosya yazılamadı", http.StatusInternalServerError)
		return
	}

	// 6. Veritabanına (attachments tablosuna) kaydet
	query := `INSERT INTO attachments (task_id, file_name, file_path) VALUES ($1, $2, $3)`
	_, err = database.DB.Exec(query, taskID, handler.Filename, filePath)
	if err != nil {
		log.Printf("DB Hatası: %v", err)
		http.Error(w, "Veritabanı kaydı başarısız", http.StatusInternalServerError)
		return
	}

	// 7. Başarılı Yanıt (JSON formatında)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Dosya başarıyla yüklendi!",
		"file_name": handler.Filename, // Orijinal ismi kullanıcıya gösterelim
		"stored_as": newFileName,      // Sunucudaki UUID'li hali
		"task_id":   taskID,
	})
}

func GetAttachmentsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. URL'den task_id'yi alıyoruz (Örn: /tasks/attachments?task_id=1)
	taskIDStr := r.URL.Query().Get("task_id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Geçersiz Task ID", http.StatusBadRequest)
		return
	}

	// 2. DB'den bu task_id'ye ait tüm ekleri çek
	query := `SELECT id, task_id, file_name, file_path, created_at FROM attachments WHERE task_id = $1`
	rows, err := database.DB.Query(query, taskID)
	if err != nil {
		log.Printf("DB Hatası: %v", err)
		http.Error(w, "Dosyalar getirilemedi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var attachments []models.Attachment
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.ID, &a.TaskID, &a.FileName, &a.FilePath, &a.CreatedAt); err != nil {
			log.Printf("Scan Hatası: %v", err)
			continue
		}
		attachments = append(attachments, a)
	}

	// 3. Sonucu JSON olarak dön
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attachments)
}
