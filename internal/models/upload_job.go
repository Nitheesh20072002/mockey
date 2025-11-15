package models

import (
	"time"
)

// UploadJob represents a file import job.
type UploadJob struct {
	ID            int       `db:"id" json:"id"`
	FileName      string    `db:"file_name" json:"file_name"`
	Status        string    `db:"status" json:"status"`
	TotalRows     int       `db:"total_rows" json:"total_rows"`
	ProcessedRows int       `db:"processed_rows" json:"processed_rows"`
	Errors        string    `db:"errors" json:"errors"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
