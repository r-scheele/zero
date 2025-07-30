package types

import "time"

// Resource represents a file or link attached to a note
type Resource struct {
	Type          string    `json:"type"`           // "file", "youtube", "url", "image", "pdf", "doc", "video"
	Name          string    `json:"name"`           // Display name
	URL           string    `json:"url"`            // File path or external URL
	Size          int64     `json:"size"`           // File size in bytes (0 for external links)
	MimeType      string    `json:"mime_type"`      // MIME type for files
	Thumbnail     string    `json:"thumbnail"`      // Thumbnail path for videos/images
	Duration      int       `json:"duration"`       // Duration in seconds for videos
	ExtractedText string    `json:"extracted_text"` // Text extracted from PDFs, docs, etc.
	UploadedAt    time.Time `json:"uploaded_at"`
}