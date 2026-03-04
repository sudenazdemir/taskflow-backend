package models

import "time"

// User: Sistemdeki kullanıcılar
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "-" işareti şifrenin API'den dışarı sızmasını engeller (Güvenlik!)
	CreatedAt    time.Time `json:"created_at"`
}
