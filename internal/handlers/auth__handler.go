package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sudenazdemir/taskflow-backend/internal/database"
	"github.com/sudenazdemir/taskflow-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONError(w, "Sadece POST metodu destekleniyor", http.StatusMethodNotAllowed)
		return
	}
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendJSONError(w, "Geçersiz JSON verisi", http.StatusBadRequest)
		return
	}
	// 1. şifreyi hashlemek güvenlik
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		sendJSONError(w, "Şifre hashlenirken hata oluştu", http.StatusInternalServerError)
		return
	}
	// 2. kullanıcıyı veritabanına kaydetmek
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at"
	err = database.DB.QueryRow(query, user.Username, user.Email, string(hashedPassword)).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		sendJSONError(w, "Kullanıcı zaten mevcut olabilir veya veritabanı hatası oluştu", http.StatusConflict)
		return
	}
	// 3. başarılı kayıt durumunda kullanıcı bilgilerini döndürmek
	user.PasswordHash = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Kullanıcı başarıyla kaydedildi",
		"user":    user,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONError(w, "Sadece POST metodu destekleniyor", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendJSONError(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	// 1. Kullanıcıyı veritabanında bul
	var user models.User
	var hashedPassword string
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1`

	err := database.DB.QueryRow(query, input.Email).
		Scan(&user.ID, &user.Username, &user.Email, &hashedPassword, &user.CreatedAt)
	if err != nil {
		sendJSONError(w, "E-posta veya şifre hatalı", http.StatusUnauthorized)
		return
	}

	// 2. Şifreyi karşılaştır (bcrypt.CompareHashAndPassword)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		sendJSONError(w, "E-posta veya şifre hatalı", http.StatusUnauthorized)
		return
	}

	// 3. Başarılı yanıt (Şimdilik Token yok, sadece başarı mesajı)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Giriş başarılı! Hoş geldin " + user.Username,
		"user":    user,
	})
}
