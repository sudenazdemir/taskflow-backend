package models

import "time"

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     int       `json:"owner_id"` // Bu satırı ekle!
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
