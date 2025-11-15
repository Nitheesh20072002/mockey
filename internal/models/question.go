package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Question represents an exam question.
type Question struct {
	ID            int             `db:"id" json:"id"`
	ExamID        int             `db:"exam_id" json:"exam_id"`
	Section       string          `db:"section" json:"section"`
	Topic         string          `db:"topic" json:"topic"`
	Year          int             `db:"year" json:"year"`
	Type          string          `db:"type" json:"type"`
	QuestionText  string          `db:"question_text" json:"question_text"`
	Choices       json.RawMessage `db:"choices" json:"choices"`               // JSONB: array of choice objects
	CorrectAnswer json.RawMessage `db:"correct_answer" json:"correct_answer"` // JSONB: correct answer(s)
	Marks         float64         `db:"marks" json:"marks"`
	NegativeMarks float64         `db:"negative_marks" json:"negative_marks"`
	Resources     json.RawMessage `db:"resources" json:"resources"` // JSONB: reference materials
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}

// Scan implements sql.Scanner interface for JSONB fields
func (q *Question) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, &q)
}

// Value implements driver.Valuer interface for JSONB fields
func (q Question) Value() (driver.Value, error) {
	return json.Marshal(q)
}
