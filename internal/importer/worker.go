package importer

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mockey/internal/db"
	"github.com/mockey/internal/models"
	"github.com/mockey/internal/repo"
)

// ProcessUploadFile performs CSV import with transaction handling
func ProcessUploadFile(jobID, examID int, path string) error {
	/*
		1) take exam id and file path which contain some questions in CSV
		2) create a test entry in tests table in DB with this exam id
		3) parse csv and for each row create a question in questions table in DB and link to exam id
		4) make a row in test_questions table to link question to test table
		All within a single transaction - rollback on any error
	*/
	r := repo.NewUploadJobRepo(db.DB)

	// mark running
	job, err := r.Get(jobID)
	if err != nil {
		return err
	}
	job.Status = "running"
	r.Update(job)

	// simple CSV parse: count rows
	f, err := os.Open(path)
	if err != nil {
		job.Status = "failed"
		job.Errors = fmt.Sprintf("open error: %v", err)
		r.Update(job)
		return err
	}
	defer f.Close()

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		job.Status = "failed"
		job.Errors = fmt.Sprintf("transaction begin error: %v", err)
		_ = r.Update(job)
		return err
	}
	defer tx.Rollback() // rollback if not committed

	// Create transaction-scoped repos
	txTestsRepo := repo.NewTestsRepoWithTx(tx)
	txQuestionsRepo := repo.NewQuestionRepoWithTx(tx)
	txTestQuestionsRepo := repo.NewTestQuestionsRepoWithTx(tx)

	// Step 1: Create a test entry for this exam (within transaction)
	test := &models.Test{
		ExamID:    examID,
		Title:     fmt.Sprintf("Test for Exam %d", examID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := txTestsRepo.Create(test); err != nil {
		job.Status = "failed"
		job.Errors = fmt.Sprintf("create test error: %v", err)
		_ = r.Update(job)
		return err
	}

	// Step 2: Parse CSV and create questions
	reader := csv.NewReader(bufio.NewReader(f))

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		job.Status = "failed"
		job.Errors = fmt.Sprintf("csv header read error: %v", err)
		_ = r.Update(job)
		return err
	}

	count := 0
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			// on read error, mark job failed and return (tx will rollback)
			job.Status = "failed"
			job.Errors = fmt.Sprintf("csv read error: %v", err)
			_ = r.Update(job)
			return err
		}

		// skip empty rows
		if len(row) == 0 {
			continue
		}

		// Expect columns: section, topic, year, type, question_text, choices (JSON), correct_answer (JSON), marks, negative_marks, resources (JSON, optional)
		section := ""
		if len(row) > 0 {
			section = strings.TrimSpace(row[0])
		}
		topic := ""
		if len(row) > 1 {
			topic = strings.TrimSpace(row[1])
		}
		year := 0
		if len(row) > 2 {
			if v, err := strconv.Atoi(strings.TrimSpace(row[2])); err == nil {
				year = v
			}
		}
		qType := ""
		if len(row) > 3 {
			qType = strings.TrimSpace(row[3])
		}
		questionText := ""
		if len(row) > 4 {
			questionText = strings.TrimSpace(row[4])
		}
		if questionText == "" {
			// skip rows with empty question
			continue
		}

		choices := json.RawMessage(nil)
		if len(row) > 5 {
			choicesStr := strings.TrimSpace(row[5])
			if choicesStr != "" {
				choices = []byte(choicesStr)
			} else {
				choices = []byte("null")
			}
		} else {
			choices = []byte("null")
		}

		correctAnswer := []byte(nil)
		if len(row) > 6 {
			correctStr := strings.TrimSpace(row[6])
			if correctStr != "" {
				correctAnswer = []byte(correctStr)
			} else {
				correctAnswer = []byte("null")
			}
		} else {
			correctAnswer = []byte("null")
		}

		marks := 1.0
		if len(row) > 7 {
			if v, err := strconv.ParseFloat(strings.TrimSpace(row[7]), 64); err == nil {
				marks = v
			}
		}
		negativeMarks := 0.0
		if len(row) > 8 {
			if v, err := strconv.ParseFloat(strings.TrimSpace(row[8]), 64); err == nil {
				negativeMarks = v
			}
		}

		resources := []byte(nil)
		if len(row) > 9 {
			resourcesStr := strings.TrimSpace(row[9])
			if resourcesStr != "" {
				resources = []byte(resourcesStr)
			} else {
				resources = []byte("null")
			}
		} else {
			resources = []byte("null")
		}

		// Step 3: Create question (within transaction)
		question := &models.Question{
			ExamID:        examID,
			Section:       section,
			Topic:         topic,
			Year:          year,
			Type:          qType,
			QuestionText:  questionText,
			Choices:       choices,
			CorrectAnswer: correctAnswer,
			Marks:         marks,
			NegativeMarks: negativeMarks,
			Resources:     resources,
			CreatedAt:     time.Now(),
		}

		if err := txQuestionsRepo.Create(question); err != nil {
			job.Status = "failed"
			job.Errors = fmt.Sprintf("create question error: %v", err)
			_ = r.Update(job)
			return err
		}

		// Step 4: Link question to test (within transaction)
		testQuestion := &models.TestQuestion{
			TestID:     test.ID,
			QuestionID: question.ID,
			CreatedAt:  time.Now(),
		}

		if err := txTestQuestionsRepo.Create(testQuestion); err != nil {
			job.Status = "failed"
			job.Errors = fmt.Sprintf("link question error: %v", err)
			_ = r.Update(job)
			return err
		}

		count++
	}

	// All operations successful - commit transaction
	if err := tx.Commit(); err != nil {
		job.Status = "failed"
		job.Errors = fmt.Sprintf("transaction commit error: %v", err)
		_ = r.Update(job)
		return err
	}

	// Update job status to finished
	job.TotalRows = count
	job.ProcessedRows = count
	job.Status = "finished"
	if err := r.Update(job); err != nil {
		return err
	}

	// Optionally remove file after processing
	_ = os.Remove(path)
	return nil
}
