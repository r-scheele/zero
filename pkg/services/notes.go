package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/r-scheele/zero/config"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/ent/note"
	"github.com/r-scheele/zero/ent/notelike"
	"github.com/r-scheele/zero/ent/noterepost"
	"github.com/r-scheele/zero/ent/user"
	"github.com/r-scheele/zero/pkg/types"
	"github.com/spf13/afero"
)

// NotesService handles note-related operations
type NotesService struct {
	orm   *ent.Client
	files afero.Fs
	cache *CacheClient
	config *config.Config
}

// NewNotesService creates a new notes service
func NewNotesService(orm *ent.Client, files afero.Fs, cache *CacheClient, config *config.Config) *NotesService {
	return &NotesService{
		orm:    orm,
		files:  files,
		cache:  cache,
		config: config,
	}
}

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

// CreateNoteInput represents input for creating a note
type CreateNoteInput struct {
	Title           string     `json:"title" validate:"required,min=1,max=200"`
	Description     string     `json:"description" validate:"max=500"`
	Content         string     `json:"content"`
	Visibility      string     `json:"visibility" validate:"oneof=private public"`
	PermissionLevel string     `json:"permission_level" validate:"oneof=read_only read_write read_write_approval"`
	Resources       []Resource `json:"resources"`
}

// UpdateNoteInput represents input for updating a note
type UpdateNoteInput struct {
	Title           *string     `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description     *string     `json:"description,omitempty" validate:"omitempty,max=500"`
	Content         *string     `json:"content,omitempty"`
	Visibility      *string     `json:"visibility,omitempty" validate:"omitempty,oneof=private public"`
	PermissionLevel *string     `json:"permission_level,omitempty" validate:"omitempty,oneof=read_only read_write read_write_approval"`
	Resources       *[]Resource `json:"resources,omitempty"`
}

// CreateNote creates a new note
func (s *NotesService) CreateNote(ctx context.Context, userID int, input CreateNoteInput) (*ent.Note, error) {
	// Generate share token
	shareToken, err := s.generateShareToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate share token: %w", err)
	}

	// Create note
	noteBuilder := s.orm.Note.Create().
		SetTitle(input.Title).
		SetVisibility(note.Visibility(input.Visibility)).
		SetPermissionLevel(note.PermissionLevel(input.PermissionLevel)).
		SetShareToken(shareToken).
		SetOwnerID(userID)

	if input.Description != "" {
		noteBuilder = noteBuilder.SetDescription(input.Description)
	}

	if input.Content != "" {
		noteBuilder = noteBuilder.SetContent(input.Content)
	}

	// Convert resources to types.Resource for database storage
	if len(input.Resources) > 0 {
		resourceData := make([]types.Resource, len(input.Resources))
		for i, r := range input.Resources {
			resourceData[i] = types.Resource{
				Type:          r.Type,
				Name:          r.Name,
				URL:           r.URL,
				Size:          r.Size,
				MimeType:      r.MimeType,
				Thumbnail:     r.Thumbnail,
				Duration:      r.Duration,
				ExtractedText: r.ExtractedText,
				UploadedAt:    r.UploadedAt,
			}
		}
		noteBuilder = noteBuilder.SetResources(resourceData)
	}

	createdNote, err := noteBuilder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	// TODO: Add AI processing later
	// For now, we'll set ai_processing to false since we don't have AI integration yet

	return createdNote, nil
}

// UpdateNote updates an existing note
func (s *NotesService) UpdateNote(ctx context.Context, noteID, userID int, input UpdateNoteInput) (*ent.Note, error) {
	// Check if user owns the note
	exists, err := s.orm.Note.Query().
		Where(note.ID(noteID), note.HasOwnerWith(user.ID(userID))).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check note ownership: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("note not found or access denied")
	}

	// Build update query
	updateBuilder := s.orm.Note.UpdateOneID(noteID)

	if input.Title != nil {
		updateBuilder = updateBuilder.SetTitle(*input.Title)
	}

	if input.Description != nil {
		updateBuilder = updateBuilder.SetDescription(*input.Description)
	}

	if input.Content != nil {
		updateBuilder = updateBuilder.SetContent(*input.Content)
	}

	if input.Visibility != nil {
		updateBuilder = updateBuilder.SetVisibility(note.Visibility(*input.Visibility))
	}

	if input.PermissionLevel != nil {
		updateBuilder = updateBuilder.SetPermissionLevel(note.PermissionLevel(*input.PermissionLevel))
	}

	// Handle resources update
	if input.Resources != nil {
		resourceData := make([]types.Resource, len(*input.Resources))
		for i, r := range *input.Resources {
			resourceData[i] = types.Resource{
				Type:          r.Type,
				Name:          r.Name,
				URL:           r.URL,
				Size:          r.Size,
				MimeType:      r.MimeType,
				Thumbnail:     r.Thumbnail,
				Duration:      r.Duration,
				ExtractedText: r.ExtractedText,
				UploadedAt:    r.UploadedAt,
			}
		}
		updateBuilder = updateBuilder.SetResources(resourceData)
	}

	updatedNote, err := updateBuilder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	// TODO: Add AI reprocessing later

	return updatedNote, nil
}

// GetNote retrieves a note by ID with permission checking
func (s *NotesService) GetNote(ctx context.Context, noteID int, userID *int) (*ent.Note, error) {
	noteQuery := s.orm.Note.Query().Where(note.ID(noteID))

	// If user is not provided, only allow public notes
	if userID == nil {
		noteQuery = noteQuery.Where(note.VisibilityEQ(note.VisibilityPublic))
	} else {
		// Allow if user owns the note or if it's public
		noteQuery = noteQuery.Where(
			note.Or(
				note.HasOwnerWith(user.ID(*userID)),
				note.VisibilityEQ(note.VisibilityPublic),
			),
		)
	}

	fetchedNote, err := noteQuery.WithOwner().Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("note not found or access denied")
		}
		return nil, fmt.Errorf("failed to fetch note: %w", err)
	}

	return fetchedNote, nil
}

// GetNoteByShareToken retrieves a note by its share token
func (s *NotesService) GetNoteByShareToken(ctx context.Context, shareToken string) (*ent.Note, error) {
	fetchedNote, err := s.orm.Note.Query().
		Where(note.ShareToken(shareToken)).
		WithOwner().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("failed to fetch note: %w", err)
	}

	return fetchedNote, nil
}

// ListUserNotes lists notes owned by a user
func (s *NotesService) ListUserNotes(ctx context.Context, userID int, limit, offset int) ([]*ent.Note, error) {
	notes, err := s.orm.Note.Query().
		Where(note.HasOwnerWith(user.ID(userID))).
		Order(ent.Desc(note.FieldUpdatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user notes: %w", err)
	}

	return notes, nil
}

// ListPublicNotes lists public notes for the home page
func (s *NotesService) ListPublicNotes(ctx context.Context, limit, offset int) ([]*ent.Note, error) {
	notes, err := s.orm.Note.Query().
		Where(note.VisibilityEQ(note.VisibilityPublic)).
		Order(ent.Desc(note.FieldUpdatedAt)).
		WithOwner().
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public notes: %w", err)
	}

	return notes, nil
}

// DeleteNote deletes a note (only by owner)
func (s *NotesService) DeleteNote(ctx context.Context, noteID, userID int) error {
	// Check if user owns the note
	exists, err := s.orm.Note.Query().
		Where(note.ID(noteID), note.HasOwnerWith(user.ID(userID))).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to check note ownership: %w", err)
	}
	if !exists {
		return fmt.Errorf("note not found or access denied")
	}

	// Delete the note
	err = s.orm.Note.DeleteOneID(noteID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}

// UploadResource handles file upload for a note
func (s *NotesService) UploadResource(ctx context.Context, noteID, userID int, file *multipart.FileHeader) (*Resource, error) {
	// Check if user owns the note
	exists, err := s.orm.Note.Query().
		Where(note.ID(noteID), note.HasOwnerWith(user.ID(userID))).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check note ownership: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("note not found or access denied")
	}

	// Validate file size (40MB max per file)
	if file.Size > 40*1024*1024 {
		return nil, fmt.Errorf("file size exceeds 40MB limit")
	}

	// Create uploads directory
	uploadsDir := fmt.Sprintf("uploads/notes/%d", noteID)
	if err := s.files.MkdirAll(uploadsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), generateRandomString(8), ext)
	filePath := filepath.Join(uploadsDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := s.files.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file
	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Determine resource type and extract text if applicable
	resourceType := s.determineResourceType(file.Header.Get("Content-Type"), ext)
	extractedText := ""

	if resourceType == "pdf" || resourceType == "doc" {
		// TODO: Implement text extraction from PDFs and docs
		// This would require additional libraries like pdfcpu or unioffice
	}

	// Create resource object
	resource := &Resource{
		Type:          resourceType,
		Name:          file.Filename,
		URL:           filePath,
		Size:          file.Size,
		MimeType:      file.Header.Get("Content-Type"),
		ExtractedText: extractedText,
		UploadedAt:    time.Now(),
	}

	return resource, nil
}

// AddURLResource adds a URL resource to a note
func (s *NotesService) AddURLResource(ctx context.Context, noteID, userID int, resourceURL, name string) (*Resource, error) {
	// Check if user owns the note
	exists, err := s.orm.Note.Query().
		Where(note.ID(noteID), note.HasOwnerWith(user.ID(userID))).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check note ownership: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("note not found or access denied")
	}

	// Validate URL
	parsedURL, err := url.Parse(resourceURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, fmt.Errorf("invalid URL")
	}

	// Determine resource type
	resourceType := "url"
	if s.isYouTubeURL(resourceURL) {
		resourceType = "youtube"
	}

	// Create resource object
	resource := &Resource{
		Type:       resourceType,
		Name:       name,
		URL:        resourceURL,
		Size:       0,
		UploadedAt: time.Now(),
	}

	return resource, nil
}

// AddResourceToNote adds a resource to an existing note
func (s *NotesService) AddResourceToNote(ctx context.Context, noteID, userID int, resource types.Resource) error {
	// Check if user owns the note
	exists, err := s.orm.Note.Query().
		Where(note.ID(noteID), note.HasOwnerWith(user.ID(userID))).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to check note ownership: %w", err)
	}
	if !exists {
		return fmt.Errorf("note not found or access denied")
	}

	// Get current note to append to existing resources
	currentNote, err := s.orm.Note.Query().
		Where(note.ID(noteID)).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch current note: %w", err)
	}

	// Calculate current total size
	currentTotalSize := s.calculateTotalResourceSize(currentNote.Resources)
	
	// Check if adding this resource would exceed 400MB limit
	if currentTotalSize+resource.Size > 400*1024*1024 {
		return fmt.Errorf("adding this resource would exceed the 400MB total limit for this note")
	}

	// Append new resource to existing resources
	existingResources := currentNote.Resources
	updatedResources := append(existingResources, resource)

	// Update note with new resources
	_, err = s.orm.Note.UpdateOneID(noteID).
		SetResources(updatedResources).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update note resources: %w", err)
	}

	return nil
}

// calculateTotalResourceSize calculates the total size of all resources
func (s *NotesService) calculateTotalResourceSize(resources []types.Resource) int64 {
	var totalSize int64
	for _, resource := range resources {
		totalSize += resource.Size
	}
	return totalSize
}

// LikeNote adds a like to a note
func (s *NotesService) LikeNote(ctx context.Context, noteID, userID int) error {
	// Check if note exists and is public
	noteExists, err := s.orm.Note.Query().
		Where(note.ID(noteID), note.VisibilityEQ(note.VisibilityPublic)).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to check note existence: %w", err)
	}
	if !noteExists {
		return fmt.Errorf("note not found or not public")
	}

	// Check if user already liked this note
	likeExists, err := s.orm.NoteLike.Query().
		Where(
			notelike.HasUserWith(user.ID(userID)),
			notelike.HasNoteWith(note.ID(noteID)),
		).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to check existing like: %w", err)
	}
	if likeExists {
		return fmt.Errorf("note already liked by user")
	}

	// Create the like
	_, err = s.orm.NoteLike.Create().
		SetUserID(userID).
		SetNoteID(noteID).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}

	// Invalidate cache for like count
	cacheKey := fmt.Sprintf("note_likes_count:%d", noteID)
	s.cache.Flush().Key(cacheKey).Execute(ctx)

	return nil
}

// UnlikeNote removes a like from a note
func (s *NotesService) UnlikeNote(ctx context.Context, noteID, userID int) error {
	// Delete the like
	_, err := s.orm.NoteLike.Delete().
		Where(
			notelike.HasUserWith(user.ID(userID)),
			notelike.HasNoteWith(note.ID(noteID)),
		).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	// Invalidate cache for like count
	cacheKey := fmt.Sprintf("note_likes_count:%d", noteID)
	s.cache.Flush().Key(cacheKey).Execute(ctx)

	return nil
}

// RepostNote creates a repost of a note
func (s *NotesService) RepostNote(ctx context.Context, noteID, userID int, comment string) error {
	// Check if note exists and is public
	noteExists, err := s.orm.Note.Query().
		Where(note.ID(noteID), note.VisibilityEQ(note.VisibilityPublic)).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to check note existence: %w", err)
	}
	if !noteExists {
		return fmt.Errorf("note not found or not public")
	}

	// Check if user already reposted this note
	repostExists, err := s.orm.NoteRepost.Query().
		Where(
			noterepost.HasUserWith(user.ID(userID)),
			noterepost.HasNoteWith(note.ID(noteID)),
		).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to check existing repost: %w", err)
	}
	if repostExists {
		return fmt.Errorf("note already reposted by user")
	}

	// Create the repost
	repostBuilder := s.orm.NoteRepost.Create().
		SetUserID(userID).
		SetNoteID(noteID)
	
	if comment != "" {
		repostBuilder = repostBuilder.SetComment(comment)
	}
	
	_, err = repostBuilder.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create repost: %w", err)
	}

	// Invalidate cache for repost count
	cacheKey := fmt.Sprintf("note_reposts_count:%d", noteID)
	s.cache.Flush().Key(cacheKey).Execute(ctx)

	return nil
}

// UnrepostNote removes a repost
func (s *NotesService) UnrepostNote(ctx context.Context, noteID, userID int) error {
	// Find and delete the repost
	deleted, err := s.orm.NoteRepost.Delete().
		Where(
			noterepost.HasUserWith(user.ID(userID)),
			noterepost.HasNoteWith(note.ID(noteID)),
		).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove repost: %w", err)
	}
	if deleted == 0 {
		return fmt.Errorf("repost not found")
	}

	// Invalidate cache for repost count
	cacheKey := fmt.Sprintf("note_reposts_count:%d", noteID)
	s.cache.Flush().Key(cacheKey).Execute(ctx)

	return nil
}

// GetNoteLikesCount returns the number of likes for a note
func (s *NotesService) GetNoteLikesCount(ctx context.Context, noteID int) (int, error) {
	// Try to get count from cache first
	cacheKey := fmt.Sprintf("note_likes_count:%d", noteID)
	if cachedCount, err := s.cache.Get().Key(cacheKey).Fetch(ctx); err == nil {
		if count, ok := cachedCount.(int); ok {
			return count, nil
		}
	}

	// If not in cache, fetch from database
	count, err := s.orm.NoteLike.Query().
		Where(notelike.HasNoteWith(note.ID(noteID))).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count likes: %w", err)
	}

	// Cache the count for future requests
	s.cache.Set().
		Key(cacheKey).
		Data(count).
		Expiration(s.config.Cache.Expiration.NoteCounts).
		Save(ctx)

	return count, nil
}

// GetNoteRepostsCount returns the number of reposts for a note
func (s *NotesService) GetNoteRepostsCount(ctx context.Context, noteID int) (int, error) {
	// Try to get count from cache first
	cacheKey := fmt.Sprintf("note_reposts_count:%d", noteID)
	if cachedCount, err := s.cache.Get().Key(cacheKey).Fetch(ctx); err == nil {
		if count, ok := cachedCount.(int); ok {
			return count, nil
		}
	}

	// If not in cache, fetch from database
	count, err := s.orm.NoteRepost.Query().
		Where(noterepost.HasNoteWith(note.ID(noteID))).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count reposts: %w", err)
	}

	// Cache the count for future requests
	s.cache.Set().
		Key(cacheKey).
		Data(count).
		Expiration(s.config.Cache.Expiration.NoteCounts).
		Save(ctx)

	return count, nil
}

// IsNoteLikedByUser checks if a note is liked by a specific user
func (s *NotesService) IsNoteLikedByUser(ctx context.Context, noteID, userID int) (bool, error) {
	exists, err := s.orm.NoteLike.Query().
		Where(
			notelike.HasUserWith(user.ID(userID)),
			notelike.HasNoteWith(note.ID(noteID)),
		).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check like status: %w", err)
	}
	return exists, nil
}

// IsNoteRepostedByUser checks if a note is reposted by a specific user
func (s *NotesService) IsNoteRepostedByUser(ctx context.Context, noteID, userID int) (bool, error) {
	exists, err := s.orm.NoteRepost.Query().
		Where(
			noterepost.HasUserWith(user.ID(userID)),
			noterepost.HasNoteWith(note.ID(noteID)),
		).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check repost status: %w", err)
	}
	return exists, nil
}

// TODO: Add AI processing functions later when AI service is implemented

// generateShareToken generates a unique share token
func (s *NotesService) generateShareToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// determineResourceType determines the type of resource based on MIME type and extension
func (s *NotesService) determineResourceType(mimeType, ext string) string {
	ext = strings.ToLower(ext)
	mimeType = strings.ToLower(mimeType)

	switch {
	case strings.Contains(mimeType, "image") || ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp":
		return "image"
	case strings.Contains(mimeType, "video") || ext == ".mp4" || ext == ".avi" || ext == ".mov" || ext == ".wmv":
		return "video"
	case strings.Contains(mimeType, "pdf") || ext == ".pdf":
		return "pdf"
	case strings.Contains(mimeType, "document") || ext == ".doc" || ext == ".docx" || ext == ".odt":
		return "doc"
	default:
		return "file"
	}
}

// isYouTubeURL checks if a URL is a YouTube URL
func (s *NotesService) isYouTubeURL(urlStr string) bool {
	youtubeRegex := regexp.MustCompile(`^https?://(www\.)?(youtube\.com|youtu\.be)/`)
	return youtubeRegex.MatchString(urlStr)
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}