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

	"github.com/mockey/exam-api/internal/db"
	"github.com/mockey/exam-api/internal/models"
	"github.com/mockey/exam-api/internal/server"
)

func setupTestDBForAdmin(t *testing.T) *gorm.DB {
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := gdb.AutoMigrate(&models.User{}, &models.Exam{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	// assign to package db.DB so handlers use it
	db.DB = gdb

	return gdb
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
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d, body: %s", w.Result().StatusCode, w.Body.String())
	}
}
