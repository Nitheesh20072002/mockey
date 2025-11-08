package repo

import (
	"github.com/mockey/exam-api/internal/models"
	"gorm.io/gorm"
)

// ExamRepo provides DB operations for exams.
type ExamRepo struct {
	db *gorm.DB
}

func NewExamRepo(db *gorm.DB) *ExamRepo {
	return &ExamRepo{db: db}
}

func (r *ExamRepo) Create(e *models.Exam) error {
	return r.db.Create(e).Error
}

func (r *ExamRepo) ListByCreator(creatorID uint) ([]models.Exam, error) {
	var out []models.Exam
	if err := r.db.Where("created_by = ?", creatorID).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
