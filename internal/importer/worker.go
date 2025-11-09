package importer

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/mockey/internal/db"
	"github.com/mockey/internal/repo"
)

// ProcessUploadFile performs a simple CSV row count and updates UploadJob status.
func ProcessUploadFile(jobID uint, path string) error {
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

	reader := csv.NewReader(bufio.NewReader(f))
	count := 0
	for {
		_, err := reader.Read()
		if err != nil {
			break
		}
		count++
		// simulate processing time
		time.Sleep(5 * time.Millisecond)
	}

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
