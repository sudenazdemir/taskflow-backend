package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	AssignedTo  int       `json:"assigned_to"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
}
