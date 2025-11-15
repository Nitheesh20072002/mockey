package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/mockey/internal/models"
)

// UserRepo provides DB operations for users.
type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(u *models.User) error {
	query := `INSERT INTO users (name, email, phone, password, role, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRow(query, u.Name, u.Email, u.Phone, u.Password, u.Role, u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
}

func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Get(&u, "SELECT id, name, email, phone, password, role, created_at, updated_at FROM users WHERE email = $1", email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(id interface{}) (*models.User, error) {
	var u models.User
	err := r.db.Get(&u, "SELECT id, name, email, phone, password, role, created_at, updated_at FROM users WHERE id = $1", id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
