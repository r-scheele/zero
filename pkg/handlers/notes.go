package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/pkg/context"
	"github.com/r-scheele/zero/pkg/form"
	"github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/msg"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/tasks"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Notes struct {
	notesService *services.NotesService
	container    *services.Container
}

func init() {
	Register(new(Notes))
}

func (h *Notes) Init(c *services.Container) error {
	h.notesService = c.Notes
	h.container = c
	return nil
}

func (h *Notes) Routes(g *echo.Group) {
	// Notes routes (require authentication and verification)
	notes := g.Group("/notes", middleware.RequireAuthentication, middleware.RequireVerification)

	// List notes
	notes.GET("", h.ListNotes).Name = routenames.Notes

	// Create note
	notes.GET("/create", h.CreateNotePage).Name = routenames.Notes + ".create"
	notes.POST("/create", h.CreateNoteSubmit).Name = routenames.Notes + ".create"

	// View note
	notes.GET("/:id", h.ViewNote).Name = routenames.Notes + ".view"

	// Edit note
	notes.GET("/:id/edit", h.EditNotePage).Name = routenames.Notes + ".edit"
	notes.POST("/:id/edit", h.EditNoteSubmit).Name = routenames.Notes + ".edit"

	// Delete note
	notes.POST("/:id/delete", h.DeleteNote).Name = routenames.Notes + ".delete"

	// Like/Unlike note
	notes.POST("/:id/like", h.LikeNote).Name = routenames.Notes + ".like"
	notes.POST("/:id/unlike", h.UnlikeNote).Name = routenames.Notes + ".unlike"

	// Repost/Unrepost note
	notes.POST("/:id/repost", h.RepostNote).Name = routenames.Notes + ".repost"
	notes.POST("/:id/unrepost", h.UnrepostNote).Name = routenames.Notes + ".unrepost"

	// Share note (public access)
	g.GET("/share/:token", h.ViewSharedNote).Name = routenames.Notes + ".share"
}

// ListNotes displays the user's notes
func (h *Notes) ListNotes(ctx echo.Context) error {
	// Get authenticated user
	user := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	userID := user.ID

	// Get pagination parameters
	page := 1
	if p := ctx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 20
	offset := (page - 1) * limit

	// Fetch user's notes
	notes, err := h.notesService.ListUserNotes(ctx.Request().Context(), userID, limit, offset)
	if err != nil {
		return fail(err, "failed to fetch notes")
	}

	return pages.ListNotes(ctx, notes, page)
}

// CreateNotePage displays the create note form
func (h *Notes) CreateNotePage(ctx echo.Context) error {
	// Get the form with configuration values
	createForm := form.Get[forms.CreateNote](ctx)

	// Set file upload limits from configuration
	maxFileSize, maxTotalSize, maxFiles, err := services.GetFileUploadLimits(h.container.Config)
	if err != nil {
		// Use default values if configuration parsing fails
		maxFileSize = 40 * 1024 * 1024   // 40MB
		maxTotalSize = 400 * 1024 * 1024 // 400MB
		maxFiles = 20
	}

	createForm.MaxFileSize = maxFileSize
	createForm.MaxTotalSize = maxTotalSize
	createForm.MaxFiles = maxFiles

	return pages.CreateNote(ctx, createForm)
}

// CreateNoteSubmit handles note creation
func (h *Notes) CreateNoteSubmit(ctx echo.Context) error {
	// Get authenticated user
	user := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	userID := user.ID
	var input forms.CreateNote

	err := form.Submit(ctx, &input)
	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.CreateNotePage(ctx)
	default:
		return err
	}

	// Create note input
	createInput := services.CreateNoteInput{
		Title:           strings.TrimSpace(input.Title),
		Description:     strings.TrimSpace(input.Description),
		Content:         strings.TrimSpace(input.Content),
		Visibility:      input.Visibility,
		PermissionLevel: input.PermissionLevel,
		Resources:       []services.Resource{},
	}

	// Set defaults
	if createInput.Visibility == "" {
		createInput.Visibility = "private"
	}
	if createInput.PermissionLevel == "" {
		createInput.PermissionLevel = "read_only"
	}

	// Create the note first
	createdNote, err := h.notesService.CreateNote(ctx.Request().Context(), userID, createInput)
	if err != nil {
		msg.Error(ctx, "Failed to create note: "+err.Error())
		return h.CreateNotePage(ctx)
	}

	// Process file uploads using task queue for better performance
	multipartForm, err := ctx.MultipartForm()
	if err == nil && multipartForm.File["files"] != nil {
		// Get file upload limits from configuration
		maxFileSize, maxTotalSize, maxFiles, err := services.GetFileUploadLimits(h.container.Config)
		if err != nil {
			msg.Error(ctx, "Configuration error: "+err.Error())
			return h.CreateNotePage(ctx)
		}

		// Check file count limit
		if len(multipartForm.File["files"]) > maxFiles {
			msg.Error(ctx, fmt.Sprintf("Maximum %d files allowed", maxFiles))
			return h.CreateNotePage(ctx)
		}

		// Calculate total size of all files
		var totalSize int64
		for _, fileHeader := range multipartForm.File["files"] {
			if fileHeader.Size > maxFileSize {
				msg.Error(ctx, fmt.Sprintf("File %s exceeds %s limit", fileHeader.Filename, services.FormatFileSize(maxFileSize)))
				return h.CreateNotePage(ctx)
			}
			totalSize += fileHeader.Size
		}

		// Check total size limit
		if totalSize > maxTotalSize {
			msg.Error(ctx, fmt.Sprintf("Total file size exceeds %s limit", services.FormatFileSize(maxTotalSize)))
			return h.CreateNotePage(ctx)
		}

		for _, fileHeader := range multipartForm.File["files"] {
			if fileHeader.Size > 0 { // Only process non-empty files
				// Save file to temporary location
				src, err := fileHeader.Open()
				if err != nil {
					msg.Error(ctx, "Failed to open file "+fileHeader.Filename+": "+err.Error())
					continue
				}
				defer src.Close()

				// Create temporary file
				tempFile, err := os.CreateTemp("", "upload_*_"+filepath.Base(fileHeader.Filename))
				if err != nil {
					msg.Error(ctx, "Failed to create temporary file for "+fileHeader.Filename+": "+err.Error())
					continue
				}
				tempPath := tempFile.Name()

				// Copy file content to temporary file
				_, err = io.Copy(tempFile, src)
				tempFile.Close()
				if err != nil {
					os.Remove(tempPath) // Clean up on error
					msg.Error(ctx, "Failed to save temporary file for "+fileHeader.Filename+": "+err.Error())
					continue
				}

				// Create file upload task for background processing
				fileUploadTask := tasks.FileUploadTask{
					NoteID:   createdNote.ID,
					UserID:   userID,
					FileName: fileHeader.Filename,
					FileSize: fileHeader.Size,
					MimeType: fileHeader.Header.Get("Content-Type"),
					TempPath: tempPath,
				}

				// Queue the file upload task
				h.container.Tasks.Add(fileUploadTask)
			}
		}
	}

	// Process URL resources
	for _, resourceURL := range input.ResourceURLs {
		if strings.TrimSpace(resourceURL) != "" {
			cleanURL := strings.TrimSpace(resourceURL)
			_, err := h.notesService.AddURLResource(ctx.Request().Context(), createdNote.ID, userID, cleanURL, "Web Resource")
			if err != nil {
				// Log error but don't fail the entire operation
				msg.Error(ctx, "Failed to add URL resource: "+err.Error())
			}
		}
	}

	form.Clear(ctx)

	// Set success message and redirect to the created note
	msg.Success(ctx, "Note created successfully!")
	return ctx.Redirect(302, ctx.Echo().Reverse(routenames.Notes+".view", createdNote.ID))
}

// ViewNote displays a specific note
func (h *Notes) ViewNote(ctx echo.Context) error {
	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	// Get user ID if authenticated
	var userID *int
	if u := ctx.Get(context.AuthenticatedUserKey); u != nil {
		if user, ok := u.(*ent.User); ok {
			userID = &user.ID
		}
	}

	// Fetch the note
	note, err := h.notesService.GetNote(ctx.Request().Context(), noteID, userID)
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	return pages.ViewNote(ctx, note)
}

// LikeNote handles liking a note
func (h *Notes) LikeNote(ctx echo.Context) error {
	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	// Get authenticated user
	user := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	userID := user.ID

	// Like the note
	err = h.notesService.LikeNote(ctx.Request().Context(), noteID, userID)
	if err != nil {
		msg.Error(ctx, "Failed to like note: "+err.Error())
	} else {
		msg.Success(ctx, "Note liked!")
	}

	// Redirect back to the note
	return ctx.Redirect(302, ctx.Echo().Reverse(routenames.Notes+".view", noteID))
}

// UnlikeNote handles unliking a note
func (h *Notes) UnlikeNote(ctx echo.Context) error {
	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	// Get authenticated user
	user := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	userID := user.ID

	// Unlike the note
	err = h.notesService.UnlikeNote(ctx.Request().Context(), noteID, userID)
	if err != nil {
		msg.Error(ctx, "Failed to unlike note: "+err.Error())
	} else {
		msg.Success(ctx, "Note unliked!")
	}

	// Redirect back to the note
	return ctx.Redirect(302, ctx.Echo().Reverse(routenames.Notes+".view", noteID))
}

// RepostNote handles reposting a note
func (h *Notes) RepostNote(ctx echo.Context) error {
	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	// Get authenticated user
	user := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	userID := user.ID

	// Get optional comment from form
	comment := strings.TrimSpace(ctx.FormValue("comment"))

	// Repost the note
	err = h.notesService.RepostNote(ctx.Request().Context(), noteID, userID, comment)
	if err != nil {
		msg.Error(ctx, "Failed to repost note: "+err.Error())
	} else {
		msg.Success(ctx, "Note reposted!")
	}

	// Redirect back to the note
	return ctx.Redirect(302, ctx.Echo().Reverse(routenames.Notes+".view", noteID))
}

// UnrepostNote handles unreposting a note
func (h *Notes) UnrepostNote(ctx echo.Context) error {
	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	// Get authenticated user
	user := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	userID := user.ID

	// Unrepost the note
	err = h.notesService.UnrepostNote(ctx.Request().Context(), noteID, userID)
	if err != nil {
		msg.Error(ctx, "Failed to unrepost note: "+err.Error())
	} else {
		msg.Success(ctx, "Repost removed!")
	}

	// Redirect back to the note
	return ctx.Redirect(302, ctx.Echo().Reverse(routenames.Notes+".view", noteID))
}

// ViewSharedNote displays a note via share token
func (h *Notes) ViewSharedNote(ctx echo.Context) error {
	shareToken := ctx.Param("token")
	if shareToken == "" {
		return echo.NewHTTPError(404, "Invalid share link")
	}

	// Fetch the note by share token
	note, err := h.notesService.GetNoteByShareToken(ctx.Request().Context(), shareToken)
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	return pages.ViewNote(ctx, note)
}

// EditNotePage displays the edit note form
func (h *Notes) EditNotePage(ctx echo.Context) error {
	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	// Get authenticated user
	user := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	userID := user.ID

	// Fetch the note
	note, err := h.notesService.GetNote(ctx.Request().Context(), noteID, &userID)
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	// Check if user owns the note
	if note.Edges.Owner.ID != userID {
		return echo.NewHTTPError(403, "Access denied")
	}

	// Populate form with existing data
	editForm := form.Get[forms.EditNote](ctx)
	if !editForm.IsSubmitted() {
		editForm.ID = note.ID
		editForm.Title = note.Title
		editForm.Description = note.Description
		editForm.Content = note.Content
		editForm.Visibility = string(note.Visibility)
		editForm.PermissionLevel = string(note.PermissionLevel)

		// Extract existing resource URLs
		var resourceURLs []string
		for _, resource := range note.Resources {
			if resource.Type == "url" || resource.Type == "youtube" {
				resourceURLs = append(resourceURLs, resource.URL)
			}
		}
		editForm.ResourceURLs = resourceURLs
	}

	// Set file upload limits from configuration
	maxFileSize, maxTotalSize, maxFiles, err := services.GetFileUploadLimits(h.container.Config)
	if err != nil {
		// Use default values if configuration parsing fails
		maxFileSize = 40 * 1024 * 1024   // 40MB
		maxTotalSize = 400 * 1024 * 1024 // 400MB
		maxFiles = 20
	}

	editForm.MaxFileSize = maxFileSize
	editForm.MaxTotalSize = maxTotalSize
	editForm.MaxFiles = maxFiles

	return pages.EditNote(ctx, editForm)
}

// EditNoteSubmit handles note editing
func (h *Notes) EditNoteSubmit(ctx echo.Context) error {
	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	// Get authenticated user
	user := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	userID := user.ID
	var input forms.EditNote

	err = form.Submit(ctx, &input)
	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.EditNotePage(ctx)
	default:
		return err
	}

	// Create update input
	title := strings.TrimSpace(input.Title)
	description := strings.TrimSpace(input.Description)
	content := strings.TrimSpace(input.Content)

	updateInput := services.UpdateNoteInput{
		Title:           &title,
		Description:     &description,
		Content:         &content,
		Visibility:      &input.Visibility,
		PermissionLevel: &input.PermissionLevel,
	}

	// Update the note
	updatedNote, err := h.notesService.UpdateNote(ctx.Request().Context(), noteID, userID, updateInput)
	if err != nil {
		msg.Error(ctx, "Failed to update note: "+err.Error())
		return h.EditNotePage(ctx)
	}

	// Process file uploads using task queue for better performance
	multipartForm, err := ctx.MultipartForm()
	if err == nil && multipartForm.File["files"] != nil {
		// Get file upload limits from configuration
		maxFileSize, maxTotalSize, maxFiles, err := services.GetFileUploadLimits(h.container.Config)
		if err != nil {
			// Use default values if configuration parsing fails
			maxFileSize = 40 * 1024 * 1024   // 40MB
			maxTotalSize = 400 * 1024 * 1024 // 400MB
			maxFiles = 20
		}

		// Validate file count
		if len(multipartForm.File["files"]) > maxFiles {
			msg.Error(ctx, fmt.Sprintf("Maximum %d files allowed", maxFiles))
			return h.EditNotePage(ctx)
		}

		// Validate total file size
		var totalSize int64
		for _, fileHeader := range multipartForm.File["files"] {
			totalSize += fileHeader.Size
			if fileHeader.Size > maxFileSize {
				msg.Error(ctx, fmt.Sprintf("File %s exceeds maximum size limit", fileHeader.Filename))
				return h.EditNotePage(ctx)
			}
		}

		if totalSize > maxTotalSize {
			msg.Error(ctx, "Total file size exceeds maximum limit")
			return h.EditNotePage(ctx)
		}

		// Process each file
		for _, fileHeader := range multipartForm.File["files"] {
			// Create temporary file
			tempDir := os.TempDir()
			tempPath := filepath.Join(tempDir, fmt.Sprintf("upload_%d_%s", updatedNote.ID, fileHeader.Filename))

			// Open uploaded file
			src, err := fileHeader.Open()
			if err != nil {
				msg.Error(ctx, "Failed to open uploaded file "+fileHeader.Filename+": "+err.Error())
				continue
			}
			defer src.Close()

			// Create temporary file
			tempFile, err := os.Create(tempPath)
			if err != nil {
				msg.Error(ctx, "Failed to create temporary file for "+fileHeader.Filename+": "+err.Error())
				continue
			}

			// Copy file content
			_, err = io.Copy(tempFile, src)
			tempFile.Close()
			if err != nil {
				os.Remove(tempPath) // Clean up on error
				msg.Error(ctx, "Failed to save temporary file for "+fileHeader.Filename+": "+err.Error())
				continue
			}

			// Create file upload task for background processing
			fileUploadTask := tasks.FileUploadTask{
				NoteID:   updatedNote.ID,
				UserID:   userID,
				FileName: fileHeader.Filename,
				FileSize: fileHeader.Size,
				MimeType: fileHeader.Header.Get("Content-Type"),
				TempPath: tempPath,
			}

			// Queue the file upload task
			h.container.Tasks.Add(fileUploadTask)
		}
	}

	// Process URL resources
	for _, resourceURL := range input.ResourceURLs {
		if strings.TrimSpace(resourceURL) != "" {
			cleanURL := strings.TrimSpace(resourceURL)
			_, err := h.notesService.AddURLResource(ctx.Request().Context(), updatedNote.ID, userID, cleanURL, "Web Resource")
			if err != nil {
				// Log error but don't fail the entire operation
				msg.Error(ctx, "Failed to add URL resource: "+err.Error())
			}
		}
	}

	// Clear form before redirect
	form.Clear(ctx)

	// Handle HTMX redirect
	redirectURL := ctx.Echo().Reverse(routenames.Notes+".view", updatedNote.ID)
	if ctx.Request().Header.Get("HX-Request") == "true" {
		// For HTMX requests, use HX-Redirect header
		ctx.Response().Header().Set("HX-Redirect", redirectURL)
		return nil
	}

	// For regular requests, use standard redirect
	return ctx.Redirect(302, redirectURL)
}

// DeleteNote handles note deletion
func (h *Notes) DeleteNote(ctx echo.Context) error {
	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(404, "Note not found")
	}

	// Get authenticated user
	user := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	userID := user.ID

	// Delete the note
	err = h.notesService.DeleteNote(ctx.Request().Context(), noteID, userID)
	if err != nil {
		msg.Error(ctx, "Failed to delete note: "+err.Error())
		return ctx.Redirect(302, ctx.Echo().Reverse(routenames.Notes+".view", noteID))
	}

	msg.Success(ctx, "Note deleted successfully!")

	// Redirect to notes list
	return ctx.Redirect(302, ctx.Echo().Reverse(routenames.Notes))
}
