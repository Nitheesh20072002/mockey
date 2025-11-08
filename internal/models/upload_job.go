package models

import (
    "time"

    "gorm.io/gorm"
)

// UploadJob represents a file import job.
type UploadJob struct {
    ID            uint           `gorm:"primaryKey" json:"id"`
    FileName      string         `gorm:"size:1024" json:"file_name"`
    Status        string         `gorm:"size:50;default:'pending'" json:"status"`
    TotalRows     int            `json:"total_rows"`
    ProcessedRows int            `json:"processed_rows"`
    Errors        string         `gorm:"type:text" json:"errors"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
