package repo

import (
	"github.com/mockey/internal/models"
	"gorm.io/gorm"
)

// UploadJobRepo provides DB operations for upload jobs.
type UploadJobRepo struct {
	db *gorm.DB
}

func NewUploadJobRepo(db *gorm.DB) *UploadJobRepo {
	return &UploadJobRepo{db: db}
}

func (r *UploadJobRepo) Create(job *models.UploadJob) error {
	return r.db.Create(job).Error
}

func (r *UploadJobRepo) Update(job *models.UploadJob) error {
	return r.db.Save(job).Error
}

func (r *UploadJobRepo) Get(id uint) (*models.UploadJob, error) {
	var j models.UploadJob
	if err := r.db.First(&j, id).Error; err != nil {
		return nil, err
	}
	return &j, nil
}
