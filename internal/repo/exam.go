package repo

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mockey/internal/models"
)

// ExamRepo provides DB operations for exams.
type ExamRepo struct {
	db *sqlx.DB
}

func NewExamRepo(db *sqlx.DB) *ExamRepo {
	return &ExamRepo{db: db}
}

func (r *ExamRepo) Create(e *models.Exam) error {
	query := `INSERT INTO exams (title, description, duration_minutes, created_by, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	row := r.db.QueryRow(query, e.Title, e.Description, e.DurationMinutes, e.CreatedBy, time.Now(), time.Now())
	err := row.Scan(&e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ExamRepo) ListByCreator(creatorID int) ([]models.Exam, error) {
	var out []models.Exam
	err := r.db.Select(&out, "SELECT id, title, description, duration_minutes, created_by, created_at, updated_at FROM exams WHERE created_by = $1", creatorID)
	if err != nil {
		return nil, err
	}
	return out, nil
}
