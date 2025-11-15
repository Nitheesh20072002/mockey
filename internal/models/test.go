package models

import (
	"time"
)

// Test represents a test/quiz based on an exam.
type Test struct {
	ID        int       `db:"id" json:"id"`
	ExamID    int       `db:"exam_id" json:"exam_id"`
	Title     string    `db:"title" json:"title"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
