package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Documents struct {
	container *services.Container
}

func init() {
	Register(new(Documents))
}

func (h *Documents) Init(c *services.Container) error {
	h.container = c
	return nil
}

func (h *Documents) Routes(g *echo.Group) {
	// Documents routes (require authentication and verification)
	docs := g.Group("/documents", middleware.RequireAuthentication, middleware.RequireVerification)
	
	// List documents
	docs.GET("", h.ListDocuments).Name = "documents.list"
	
	// Upload documents
	docs.GET("/upload", h.UploadDocumentsPage).Name = "documents.upload"
	docs.POST("/upload", h.UploadDocumentsSubmit).Name = "documents.upload.submit"
	
	// View document
	docs.GET("/:id", h.ViewDocument).Name = "documents.view"
	
	// Delete document
	docs.POST("/:id/delete", h.DeleteDocument).Name = "documents.delete"
}

// ListDocuments displays all documents
func (h *Documents) ListDocuments(ctx echo.Context) error {
	return pages.DocumentsList(ctx)
}

// UploadDocumentsPage displays the document upload form
func (h *Documents) UploadDocumentsPage(ctx echo.Context) error {
	return pages.UploadDocuments(ctx)
}

// UploadDocumentsSubmit handles document upload
func (h *Documents) UploadDocumentsSubmit(ctx echo.Context) error {
	// TODO: Implement document upload logic
	return ctx.String(200, "Document upload functionality coming soon!")
}

// ViewDocument displays a specific document
func (h *Documents) ViewDocument(ctx echo.Context) error {
	return pages.ViewDocument(ctx)
}

// DeleteDocument handles document deletion
func (h *Documents) DeleteDocument(ctx echo.Context) error {
	// TODO: Implement document deletion logic
	return ctx.String(200, "Document deletion functionality coming soon!")
}