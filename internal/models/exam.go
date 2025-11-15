package models

import (
	"time"
)

// Exam represents exam metadata.
type Exam struct {
    ID               int       `db:"id" json:"id"`
    Title            string    `db:"title" json:"title"`
    Description      string    `db:"description" json:"description"`
    DurationMinutes  int       `db:"duration_minutes" json:"duration_minutes"`
    CreatedBy        int       `db:"created_by" json:"created_by"`
    CreatedAt        time.Time `db:"created_at" json:"created_at"`
    UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
