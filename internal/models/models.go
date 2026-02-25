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

// Project: Kullanıcıların oluşturduğu projeler
type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     int       `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// ProjectMember: Collaborative (İşbirlikçi) yapımızın kalbi (Many-to-Many köprüsü)
type ProjectMember struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	UserID    int       `json:"user_id"`
	Role      string    `json:"role"` // 'owner', 'editor', 'viewer'
	JoinedAt  time.Time `json:"joined_at"`
}

// Task: Projelerin içindeki görevler
type Task struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	AssignedTo  int       `json:"assigned_to"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`   // 'todo', 'in_progress', 'done'
	Priority    int       `json:"priority"` // Örn: 1 (Düşük), 2 (Orta), 3 (Yüksek)
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
}
