package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mockey/internal/db"
	"github.com/mockey/internal/models"
	"github.com/mockey/internal/repo"
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

	// extract user id set by JWT middleware (claims 'sub')
	createdBy := 0
	if uid, ok := c.Get("user_id"); ok {
		switch v := uid.(type) {
		case float64:
			createdBy = int(v)
		case int:
			createdBy = v
		case int64:
			createdBy = int(v)
		case string:
			if n, err := strconv.Atoi(v); err == nil {
				createdBy = n
			}
		}
	}

	exam := models.Exam{
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.TimeLimit,
		CreatedBy:       createdBy,
	}

	er := repo.NewExamRepo(db.DB)
	if err := er.Create(&exam); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create exam"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"exam": gin.H{"id": exam.ID, "title": exam.Title}})
}
