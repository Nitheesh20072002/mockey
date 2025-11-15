package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/mockey/internal/db"
	"github.com/mockey/internal/server"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	testDB, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Create users table for testing
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
	`

	if _, err := testDB.Exec(schema); err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	// assign to package db.DB so handlers use it
	db.DB = testDB

	return testDB
}

func TestRegisterAndLogin(t *testing.T) {
	// ensure Gin is in test mode
	gin.SetMode(gin.TestMode)

	// set JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	// setup in-memory DB
	testDB := setupTestDB(t)

	// override global DB in package
	_ = testDB

	r := gin.Default()
	server.SetupRoutes(r)

	// Register
	regBody := map[string]string{"name": "Test", "email": "t@example.com", "password": "pass123"}
	buf, _ := json.Marshal(regBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d, body: %s", w.Result().StatusCode, w.Body.String())
	}

	// Login
	loginBody := map[string]string{"email": "t@example.com", "password": "pass123"}
	buf, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on login, got %d, body: %s", w.Result().StatusCode, w.Body.String())
	}

	// parse response and ensure token present
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}
	if _, ok := resp["token"]; !ok {
		t.Fatalf("expected token in login response, body: %s", w.Body.String())
	}
}
