package models

import "time"

type Attachment struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}
