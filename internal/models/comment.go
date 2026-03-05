package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username,omitempty"` // Yorumu yapanın adını da görelim
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
