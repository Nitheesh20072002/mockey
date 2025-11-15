package models

import (
	"time"
)

// TestQuestion links a question to a test.
type TestQuestion struct {
	ID         int       `db:"id" json:"id"`
	TestID     int       `db:"test_id" json:"test_id"`
	QuestionID int       `db:"question_id" json:"question_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
