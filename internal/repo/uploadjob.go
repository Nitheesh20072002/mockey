package repo

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mockey/internal/models"
)

// UploadJobRepo provides DB operations for upload jobs.
type UploadJobRepo struct {
	db *sqlx.DB
}

func NewUploadJobRepo(db *sqlx.DB) *UploadJobRepo {
	return &UploadJobRepo{db: db}
}

func (r *UploadJobRepo) Create(job *models.UploadJob) error {
	query := `INSERT INTO upload_jobs (file_name, status, total_rows, processed_rows, errors, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRow(query, job.FileName, job.Status, job.TotalRows, job.ProcessedRows, job.Errors,time.Now(),time.Now()).Scan(&job.ID)
}

func (r *UploadJobRepo) Update(job *models.UploadJob) error {
	query := `UPDATE upload_jobs SET file_name = $1, status = $2, total_rows = $3, processed_rows = $4, errors = $5, updated_at = $6 WHERE id = $7`
	_, err := r.db.Exec(query, job.FileName, job.Status, job.TotalRows, job.ProcessedRows, job.Errors, time.Now(), job.ID)
	return err
}

func (r *UploadJobRepo) Get(id int) (*models.UploadJob, error) {
	var j models.UploadJob
	err := r.db.Get(&j, "SELECT id, file_name, status, total_rows, processed_rows, errors, created_at, updated_at FROM upload_jobs WHERE id = $1", id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &j, nil
}
