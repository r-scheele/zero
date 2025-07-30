package tasks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mikestefanello/backlite"
	"github.com/r-scheele/zero/pkg/log"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/types"
)

// FileUploadTask represents a file upload task
type FileUploadTask struct {
	NoteID   int    `json:"note_id"`
	UserID   int    `json:"user_id"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
	TempPath string `json:"temp_path"` // Temporary local path before cloud upload
}

// Config satisfies the backlite.Task interface by providing configuration for the queue
func (t FileUploadTask) Config() backlite.QueueConfig {
	return backlite.QueueConfig{
		Name:        "FileUploadTask",
		MaxAttempts: 3,
		Timeout:     5 * time.Minute,
		Backoff:     30 * time.Second,
		Retention: &backlite.Retention{
			Duration:   24 * time.Hour,
			OnlyFailed: false,
			Data: &backlite.RetainData{
				OnlyFailed: false,
			},
		},
	}
}

// NewFileUploadTaskQueue creates a new file upload task queue
func NewFileUploadTaskQueue(c *services.Container) backlite.Queue {
	return backlite.NewQueue[FileUploadTask](func(ctx context.Context, task FileUploadTask) error {
		logger := log.Default()
		logger.Info("Processing file upload task",
			"note_id", task.NoteID,
			"user_id", task.UserID,
			"file_name", task.FileName,
			"file_size", task.FileSize,
		)

		// Get services from container
		notesService := services.NewNotesService(c.ORM, c.Files, c.Cache, c.Config)

		// Upload file to cloud storage
		var fileURL string
		var err error

		// Check if we have a temporary file to upload
		if task.TempPath != "" {
			// Open the temporary file
			file, openErr := os.Open(task.TempPath)
			if openErr != nil {
				logger.Error("Failed to open temporary file", "error", openErr, "temp_path", task.TempPath)
				return fmt.Errorf("failed to open temporary file: %w", openErr)
			}
			defer file.Close()
			defer os.Remove(task.TempPath) // Clean up temp file

			// Generate unique filename for cloud storage
			uniqueFileName := generateUniqueFileName(task.FileName)
			cloudKey := fmt.Sprintf("notes/%d/%s", task.NoteID, uniqueFileName)

			// Upload to cloud storage
			fileURL, err = c.Storage.UploadFile(ctx, cloudKey, file, task.MimeType)
			if err != nil {
				logger.Error("Failed to upload file to cloud storage", "error", err, "file_name", task.FileName)
				return fmt.Errorf("failed to upload file to cloud storage: %w", err)
			}

			logger.Info("File uploaded to cloud storage", "file_name", task.FileName, "url", fileURL)
		} else {
			// Fallback: generate a placeholder URL (this shouldn't happen in normal operation)
			uniqueFileName := generateUniqueFileName(task.FileName)
			fileURL = "/files/" + uniqueFileName
			logger.Warn("No temporary file path provided, using placeholder URL", "file_name", task.FileName)
		}
		
		// Determine resource type
		resourceType := determineResourceType(task.MimeType, task.FileName)

		// Create resource record
		resource := types.Resource{
			Type:       resourceType,
			Name:       task.FileName,
			URL:        fileURL,
			Size:       task.FileSize,
			MimeType:   task.MimeType,
			UploadedAt: time.Now(),
		}

		// Add resource to note
		err = notesService.AddResourceToNote(ctx, task.NoteID, task.UserID, resource)
		if err != nil {
			logger.Error("Failed to add resource to note",
				"error", err,
				"note_id", task.NoteID,
				"file_name", task.FileName,
			)
			return fmt.Errorf("failed to add resource to note: %w", err)
		}

		logger.Info("File upload task completed successfully",
			"note_id", task.NoteID,
			"file_name", task.FileName,
			"resource_url", resource.URL,
		)

		return nil
	})
}

// generateUniqueFileName generates a unique filename to prevent conflicts
func generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

// determineResourceType determines the resource type based on MIME type and filename
func determineResourceType(mimeType, fileName string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio"
	case mimeType == "application/pdf":
		return "pdf"
	case strings.Contains(mimeType, "document") || strings.Contains(mimeType, "word"):
		return "document"
	case strings.Contains(mimeType, "text"):
		return "text"
	default:
		return "file"
	}
}