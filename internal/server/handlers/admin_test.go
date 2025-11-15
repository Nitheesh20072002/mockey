package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/mockey/internal/db"
	"github.com/mockey/internal/server"
)

func setupTestDBForAdmin(t *testing.T) *sqlx.DB {
	testDB, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Create tables for testing
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		phone TEXT,
		password TEXT NOT NULL,
		role TEXT DEFAULT 'student',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS exams (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		duration_minutes INTEGER DEFAULT 0,
		created_by INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS upload_jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_name TEXT,
		status TEXT DEFAULT 'pending',
		total_rows INTEGER,
		processed_rows INTEGER,
		errors TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := testDB.Exec(schema); err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	// assign to package db.DB so handlers use it
	db.DB = testDB

	return testDB
}

func TestCreateExam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret")

	setupTestDBForAdmin(t)

	r := gin.Default()
	server.SetupRoutes(r)

	body := map[string]interface{}{"title": "Sample Exam", "description": "Desc", "time_limit_minutes": 60}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/exams", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	// create a valid JWT for the protected admin route
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  1,
		"role": "admin",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	s, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d, body: %s", w.Result().StatusCode, w.Body.String())
	}
}
