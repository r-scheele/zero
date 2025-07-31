package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// DocumentsList displays the list of documents
func DocumentsList(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "Documents"
	r.Metatags.Description = "Browse and manage your documents"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("flex justify-between items-center mb-6"),
			H1(Class("text-3xl font-bold text-gray-900"), Text("Documents")),
			A(
				Href("/documents/upload"),
				Class("bg-emerald-600 hover:bg-emerald-700 text-white px-4 py-2 rounded-lg font-medium"),
				Text("Upload Documents"),
			),
		),
		Div(
			Class("bg-white rounded-lg shadow p-6"),
			P(Class("text-gray-600 text-center py-8"), Text("Document management functionality is coming soon! You'll be able to upload, organize, and share your documents.")),
		),
	)

	return r.Render(layouts.Primary, page)
}

// UploadDocuments displays the document upload form
func UploadDocuments(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "Upload Documents"
	r.Metatags.Description = "Upload and manage your documents"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("max-w-2xl mx-auto"),
			H1(Class("text-3xl font-bold text-gray-900 mb-6"), Text("Upload Documents")),
			Div(
				Class("bg-white rounded-lg shadow p-6"),
				Div(
					Class("text-center py-12"),
					Div(
						Class("text-6xl mb-4"),
						Text("📄"),
					),
					H2(Class("text-2xl font-semibold text-gray-900 mb-4"), Text("Document Upload Coming Soon!")),
					P(Class("text-gray-600 mb-6"), Text("We're building a comprehensive document management system that will allow you to upload, organize, and share files with advanced search and collaboration features.")),
					A(
						Href("/dashboard"),
						Class("bg-emerald-600 hover:bg-emerald-700 text-white px-6 py-2 rounded-lg font-medium"),
						Text("Back to Dashboard"),
					),
				),
			),
		),
	)

	return r.Render(layouts.Primary, page)
}

// ViewDocument displays a specific document
func ViewDocument(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "View Document"
	r.Metatags.Description = "View document details"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("bg-white rounded-lg shadow p-6"),
			P(Class("text-gray-600 text-center py-8"), Text("Document viewing functionality is coming soon!")),
		),
	)

	return r.Render(layouts.Primary, page)
}