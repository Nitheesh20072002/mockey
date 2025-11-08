package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mockey/exam-api/internal/db"
	"github.com/mockey/exam-api/internal/models"
	"github.com/mockey/exam-api/internal/repo"
	"github.com/mockey/exam-api/internal/importer"
)

// UploadQuestions accepts a multipart file and returns a job id for async processing.
func UploadQuestions(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = os.TempDir()
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload dir"})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename))
	outPath := filepath.Join(uploadDir, filename)

	out, err := os.Create(outPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	job := models.UploadJob{FileName: outPath, Status: "pending"}
	r := repo.NewUploadJobRepo(db.DB)
	if err := r.Create(&job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload job"})
		return
	}

	// start background processing (non-blocking)
	go func(j models.UploadJob) {
		_ = importer.ProcessUploadFile(j.ID, j.FileName)
	}(job)

	c.JSON(http.StatusAccepted, gin.H{"job_id": job.ID})
}
