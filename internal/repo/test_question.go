package repo

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mockey/internal/models"
)

// TestQuestionsRepo provides DB operations for test_questions links.
type TestQuestionsRepo struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewTestQuestionsRepo(db *sqlx.DB) *TestQuestionsRepo {
	return &TestQuestionsRepo{db: db}
}

func NewTestQuestionsRepoWithTx(tx *sqlx.Tx) *TestQuestionsRepo {
	return &TestQuestionsRepo{tx: tx}
}

func (r *TestQuestionsRepo) Create(tq *models.TestQuestion) error {
	if tq.CreatedAt.IsZero() {
		tq.CreatedAt = time.Now()
	}
	query := `INSERT INTO test_questions (test_id, question_id, created_at) 
	         VALUES ($1, $2, $3) RETURNING id`

	if r.tx != nil {
		return r.tx.QueryRow(query, tq.TestID, tq.QuestionID, tq.CreatedAt).Scan(&tq.ID)
	}
	return r.db.QueryRow(query, tq.TestID, tq.QuestionID, tq.CreatedAt).Scan(&tq.ID)
}

func (r *TestQuestionsRepo) GetByID(id int) (*models.TestQuestion, error) {
	var tq models.TestQuestion
	var err error

	if r.tx != nil {
		err = r.tx.Get(&tq, "SELECT id, test_id, question_id, created_at FROM test_questions WHERE id = $1", id)
	} else {
		err = r.db.Get(&tq, "SELECT id, test_id, question_id, created_at FROM test_questions WHERE id = $1", id)
	}

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &tq, nil
}

func (r *TestQuestionsRepo) ListByTest(testID int) ([]models.TestQuestion, error) {
	var tqs []models.TestQuestion
	var err error

	if r.tx != nil {
		err = r.tx.Select(&tqs, "SELECT id, test_id, question_id, created_at FROM test_questions WHERE test_id = $1 ORDER BY created_at ASC", testID)
	} else {
		err = r.db.Select(&tqs, "SELECT id, test_id, question_id, created_at FROM test_questions WHERE test_id = $1 ORDER BY created_at ASC", testID)
	}

	if err != nil {
		return nil, err
	}
	return tqs, nil
}
