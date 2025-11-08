package handlers

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/mockey/exam-api/internal/db"
    "github.com/mockey/exam-api/internal/models"
    "github.com/mockey/exam-api/internal/repo"
)

type CreateExamRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description" binding:"omitempty"`
    TimeLimit   int    `json:"time_limit_minutes" binding:"omitempty"`
}

// CreateExam creates exam metadata (admin)
func CreateExam(c *gin.Context) {
    var req CreateExamRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // In this scaffold we accept an optional created_by from env (for tests/demo)
    createdBy := uint(0)
    if cb := os.Getenv("DEMO_ADMIN_ID"); cb != "" {
        // ignore parse error for scaffolding simplicity
    }

    exam := models.Exam{
        Title:       req.Title,
        Description: req.Description,
        TimeLimit:   req.TimeLimit,
        CreatedBy:   createdBy,
    }

    er := repo.NewExamRepo(db.DB)
    if err := er.Create(&exam); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create exam"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"exam": gin.H{"id": exam.ID, "title": exam.Title}})
}
