package handlers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mockey/internal/db"
	"github.com/mockey/internal/importer"
	"github.com/mockey/internal/models"
	"github.com/mockey/internal/repo"
)

// UploadQuestions accepts a multipart file and exam_id query param, returns a job id for async processing.
func UploadQuestions(c *gin.Context) {
	// Get exam_id from query parameter
	examIDStr := c.Query("exam_id")
	if examIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exam_id query parameter is required"})
		return
	}
	examID, err := strconv.Atoi(examIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exam_id must be an integer"})
		return
	}

	// Parse file from multipart form
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Open file stream (does NOT load into memory)
	src, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer src.Close()

	// Use $UPLOAD_DIR or system temp directory
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = os.TempDir()
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload dir"})
		return
	}

	// Unique filename to avoid collisions
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename))
	outPath := filepath.Join(uploadDir, filename)

	// Create output file
	dst, err := os.Create(outPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create file"})
		return
	}
	defer dst.Close()

	// Stream the file → avoids memory spike
	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Count total rows in CSV (excluding header)
	file, err := os.Open(outPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file for row count"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read CSV header"})
		return
	}

	rowCount := 0
	for {
		_, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			// If there's an error reading, we still create the job with 0 total rows
			break
		}
		rowCount++
	}

	// Create DB job entry with total rows
	job := models.UploadJob{
		FileName:  outPath,
		Status:    "pending",
		TotalRows: rowCount,
	}
	r := repo.NewUploadJobRepo(db.DB)
	if err := r.Create(&job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload job"})
		return
	}
	err = importer.ProcessUploadFile(job.ID, examID, job.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"job_id": job.ID, "total_rows": rowCount})
}
