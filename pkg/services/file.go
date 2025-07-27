package services

import (
	"fmt"
	"io"
	"time"

	"github.com/r-scheele/zero/pkg/ui/models"
	"github.com/spf13/afero"
)

// FileService handles file operations
type FileService struct {
	fs afero.Fs
}

// NewFileService creates a new file service
func NewFileService(fs afero.Fs) *FileService {
	return &FileService{
		fs: fs,
	}
}

// ListFiles returns a list of uploaded files
func (s *FileService) ListFiles() ([]*models.File, error) {
	// Get list of uploaded files
	info, err := afero.ReadDir(s.fs, "")
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	files := make([]*models.File, 0)
	for _, file := range info {
		files = append(files, &models.File{
			Name:     file.Name(),
			Size:     file.Size(),
			Modified: file.ModTime().Format(time.DateTime),
		})
	}

	return files, nil
}

// UploadFile handles file upload
func (s *FileService) UploadFile(filename string, src io.Reader, size int64) error {
	dst, err := s.fs.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// FileExists checks if a file exists
func (s *FileService) FileExists(filename string) (bool, error) {
	exists, err := afero.Exists(s.fs, filename)
	if err != nil {
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}
	return exists, nil
}

// DeleteFile deletes a file
func (s *FileService) DeleteFile(filename string) error {
	err := s.fs.Remove(filename)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetFileInfo returns information about a file
func (s *FileService) GetFileInfo(filename string) (*models.File, error) {
	info, err := s.fs.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &models.File{
		Name:     info.Name(),
		Size:     info.Size(),
		Modified: info.ModTime().Format(time.DateTime),
	}, nil
}