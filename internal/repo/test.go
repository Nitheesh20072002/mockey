package repo

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mockey/internal/models"
)

// TestsRepo provides DB operations for tests.
type TestsRepo struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewTestsRepo(db *sqlx.DB) *TestsRepo {
	return &TestsRepo{db: db}
}

func NewTestsRepoWithTx(tx *sqlx.Tx) *TestsRepo {
	return &TestsRepo{tx: tx}
}

func (r *TestsRepo) Create(t *models.Test) error {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = time.Now()
	}
	query := `INSERT INTO tests (exam_id, title, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4) RETURNING id`

	if r.tx != nil {
		return r.tx.QueryRow(query, t.ExamID, t.Title, t.CreatedAt, t.UpdatedAt).Scan(&t.ID)
	}
	return r.db.QueryRow(query, t.ExamID, t.Title, t.CreatedAt, t.UpdatedAt).Scan(&t.ID)
}

func (r *TestsRepo) GetByID(id int) (*models.Test, error) {
	var t models.Test
	var err error

	if r.tx != nil {
		err = r.tx.Get(&t, "SELECT id, exam_id, title, created_at, updated_at FROM tests WHERE id = $1", id)
	} else {
		err = r.db.Get(&t, "SELECT id, exam_id, title, created_at, updated_at FROM tests WHERE id = $1", id)
	}

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *TestsRepo) ListByExam(examID int) ([]models.Test, error) {
	var tests []models.Test
	var err error

	if r.tx != nil {
		err = r.tx.Select(&tests, "SELECT id, exam_id, title, created_at, updated_at FROM tests WHERE exam_id = $1 ORDER BY created_at DESC", examID)
	} else {
		err = r.db.Select(&tests, "SELECT id, exam_id, title, created_at, updated_at FROM tests WHERE exam_id = $1 ORDER BY created_at DESC", examID)
	}

	if err != nil {
		return nil, err
	}
	return tests, nil
}
