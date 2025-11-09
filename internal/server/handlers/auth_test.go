package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mockey/internal/db"
	"github.com/mockey/internal/models"
	"github.com/mockey/internal/server"
)

func setupTestDB(t *testing.T) *gorm.DB {
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := gdb.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	// assign to package db.DB so handlers use it
	db.DB = gdb

	return gdb
}

func TestRegisterAndLogin(t *testing.T) {
	// ensure Gin is in test mode
	gin.SetMode(gin.TestMode)

	// set JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	// setup in-memory DB
	gdb := setupTestDB(t)

	// override global DB in package
	_ = gdb

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
