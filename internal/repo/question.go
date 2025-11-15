package repo

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mockey/internal/models"
)

// QuestionRepo provides DB operations for questions.
type QuestionRepo struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewQuestionRepo(db *sqlx.DB) *QuestionRepo {
	return &QuestionRepo{db: db}
}

func NewQuestionRepoWithTx(tx *sqlx.Tx) *QuestionRepo {
	return &QuestionRepo{tx: tx}
}

func (r *QuestionRepo) Create(q *models.Question) error {
	if q.CreatedAt.IsZero() {
		q.CreatedAt = time.Now()
	}
	query := `INSERT INTO questions (exam_id, section, topic, year, type, question_text, choices, correct_answer, marks, negative_marks, resources, created_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`

	if r.tx != nil {
		return r.tx.QueryRow(query, q.ExamID, q.Section, q.Topic, q.Year, q.Type, q.QuestionText, q.Choices, q.CorrectAnswer, q.Marks, q.NegativeMarks, q.Resources, time.Now()).Scan(&q.ID)
	}
	return r.db.QueryRow(query, q.ExamID, q.Section, q.Topic, q.Year, q.Type, q.QuestionText, q.Choices, q.CorrectAnswer, q.Marks, q.NegativeMarks, q.Resources, time.Now()).Scan(&q.ID)
}

func (r *QuestionRepo) GetByID(id int) (*models.Question, error) {
	var q models.Question
	var err error

	if r.tx != nil {
		err = r.tx.Get(&q, "SELECT id, exam_id, section, topic, year, type, question_text, choices, correct_answer, marks, negative_marks, resources, created_at FROM questions WHERE id = $1", id)
	} else {
		err = r.db.Get(&q, "SELECT id, exam_id, section, topic, year, type, question_text, choices, correct_answer, marks, negative_marks, resources, created_at FROM questions WHERE id = $1", id)
	}

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func (r *QuestionRepo) ListByExam(examID int) ([]models.Question, error) {
	var questions []models.Question
	var err error

	if r.tx != nil {
		err = r.tx.Select(&questions, "SELECT id, exam_id, section, topic, year, type, question_text, choices, correct_answer, marks, negative_marks, resources, created_at FROM questions WHERE exam_id = $1 ORDER BY created_at ASC", examID)
	} else {
		err = r.db.Select(&questions, "SELECT id, exam_id, section, topic, year, type, question_text, choices, correct_answer, marks, negative_marks, resources, created_at FROM questions WHERE exam_id = $1 ORDER BY created_at ASC", examID)
	}

	if err != nil {
		return nil, err
	}
	return questions, nil
}
