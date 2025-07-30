package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/r-scheele/zero/config"
)

// ParseFileSize converts a size string (e.g., "40MB", "1GB") to bytes
func ParseFileSize(sizeStr string) (int64, error) {
	if sizeStr == "" {
		return 0, fmt.Errorf("empty size string")
	}

	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))
	
	// Extract number and unit
	var numStr string
	var unit string
	
	for i, char := range sizeStr {
		if char >= '0' && char <= '9' || char == '.' {
			numStr += string(char)
		} else {
			unit = sizeStr[i:]
			break
		}
	}
	
	if numStr == "" {
		return 0, fmt.Errorf("invalid size format: %s", sizeStr)
	}
	
	size, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number in size: %s", numStr)
	}
	
	// Convert to bytes based on unit
	switch unit {
	case "B", "BYTES":
		return int64(size), nil
	case "KB", "K":
		return int64(size * 1024), nil
	case "MB", "M":
		return int64(size * 1024 * 1024), nil
	case "GB", "G":
		return int64(size * 1024 * 1024 * 1024), nil
	case "TB", "T":
		return int64(size * 1024 * 1024 * 1024 * 1024), nil
	default:
		return 0, fmt.Errorf("unknown unit: %s", unit)
	}
}

// GetFileUploadLimits returns the file upload limits from configuration
func GetFileUploadLimits(cfg *config.Config) (maxFileSize, maxTotalSize int64, maxFiles int, err error) {
	maxFileSize, err = ParseFileSize(cfg.FileUpload.MaxFileSize)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid maxFileSize: %w", err)
	}
	
	maxTotalSize, err = ParseFileSize(cfg.FileUpload.MaxTotalSize)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid maxTotalSize: %w", err)
	}
	
	maxFiles = cfg.FileUpload.MaxFiles
	if maxFiles <= 0 {
		maxFiles = 20 // Default fallback
	}
	
	return maxFileSize, maxTotalSize, maxFiles, nil
}

// FormatFileSize converts bytes to a human-readable string
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}