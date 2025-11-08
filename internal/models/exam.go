package models

import (
    "time"

    "gorm.io/gorm"
)

// Exam represents exam metadata.
type Exam struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Title       string         `gorm:"size:255;not null" json:"title"`
    Description string         `gorm:"type:text" json:"description"`
    TimeLimit   int            `gorm:"default:0" json:"time_limit_minutes"` // minutes
    CreatedBy   uint           `json:"created_by"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
