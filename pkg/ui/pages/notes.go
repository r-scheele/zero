package pages

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// ListNotes displays a list of user's notes
func ListNotes(ctx echo.Context, notes []*ent.Note, page int) error {
	r := ui.NewRequest(ctx)


	content := Div(
		Class("space-y-6"),

		// Header with create button (only when notes exist)
		If(len(notes) > 0,
			Div(
				Class("flex justify-end items-center mb-6"),
				A(
					Href(r.Path(routenames.Notes+".create")),
					Class("inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"),
					Text("Create Note"),
				),
			),
		),

		// Notes grid
		If(len(notes) > 0,
			Div(
				Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"),
				Group(func() []Node {
					var noteCards []Node
					for _, note := range notes {
						noteCards = append(noteCards, noteCard(r, note))
					}
					return noteCards
				}()),
			),
		),

		// Empty state
		If(len(notes) == 0,
			Div(
				Class("text-center py-12"),
				P(
					Class("text-gray-600 mb-6"),
					Text("Create your first AI powered note to get started."),
				),
				A(
					Href(r.Path(routenames.Notes+".create")),
					Class("inline-flex items-center px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"),
					Text("Create Your First Note"),
				),
			),
		),
	)

	return r.Render(layouts.Primary, content)
}

// CreateNote displays the create note form
func CreateNote(ctx echo.Context, form *forms.CreateNote) error {
	r := ui.NewRequest(ctx)
	r.Title = "Create Note"

	content := Div(
		Class("max-w-4xl mx-auto"),



		// Form
		form.Render(r),
	)

	return r.Render(layouts.Primary, content)
}

// ViewNote displays a specific note
func ViewNote(ctx echo.Context, note *ent.Note) error {
	r := ui.NewRequest(ctx)
	r.Title = note.Title

	// Check if user owns the note
	isOwner := false
	if r.AuthUser != nil && note.Edges.Owner != nil {
		isOwner = r.AuthUser.ID == note.Edges.Owner.ID
	}

	content := Div(
		Class("max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8"),

		// Header section
		Div(
			Class("mb-8"),
			H1(
				Class("text-2xl sm:text-3xl font-bold text-gray-900 mb-4 break-words"),
				Text(note.Title),
			),
			Div(
				Class("flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4"),
				Div(
					Class("flex flex-wrap items-center gap-4 text-sm text-gray-600"),
					If(note.Edges.Owner != nil,
						Span(
							Class("font-medium text-gray-900"),
							Text("By "+note.Edges.Owner.Name),
						),
					),
					Span(
						Text(note.CreatedAt.Format("Jan 2, 2006")),
					),
					If(string(note.Visibility) == "public",
						Span(
							Class("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"),
							Text("Public"),
						),
					),
					If(string(note.Visibility) == "private",
						Span(
							Class("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800"),
							Text("Private"),
						),
					),
				),
				// Actions (only for owner)
				If(isOwner,
					Div(
						Class("flex items-center gap-2"),
						// Share button
						Button(
							Type("button"),
							Class("px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"),
							Attr("onclick", fmt.Sprintf("copyShareLink('%s')", note.ShareToken)),
							Text("Share"),
						),
						// Edit button
						A(
							Href(r.Path(routenames.Notes+".edit", note.ID)),
							Class("px-3 py-1.5 text-sm font-medium text-blue-700 bg-blue-50 border border-blue-300 rounded-md hover:bg-blue-100 transition-colors"),
							Text("Edit"),
						),
						// Delete button
						Form(
							Class("inline"),
							Method("POST"),
							Action(r.Path(routenames.Notes+".delete", note.ID)),
							Attr("onsubmit", "return confirm('Are you sure you want to delete this note?')"),
							CSRF(r),
							Button(
								Type("submit"),
								Class("px-3 py-1.5 text-sm font-medium text-red-700 bg-red-50 border border-red-300 rounded-md hover:bg-red-100 transition-colors"),
								Text("Delete"),
							),
						),
					),
				),
			),
		),

		// Description
		If(note.Description != "",
			Div(
				Class("mb-6 p-4 bg-blue-50 border-l-4 border-blue-400 rounded-r-md"),
				H3(
					Class("text-sm font-semibold text-blue-900 mb-2"),
					Text("Description"),
				),
				P(
					Class("text-blue-800 leading-relaxed"),
					Text(note.Description),
				),
			),
		),

		// Content
		If(note.Content != "",
			Div(
				Class("mb-6"),
				H2(
					Class("text-lg font-semibold text-gray-900 mb-4"),
					Text("Content"),
				),
				Div(
					Class("bg-white border border-gray-200 rounded-lg p-6"),
					Pre(
						Class("whitespace-pre-wrap text-gray-700 leading-relaxed break-words font-sans text-base"),
						Text(note.Content),
					),
				),
			),
		),

		// AI Curriculum (if available)
		If(note.AiCurriculum != "",
			Div(
				Class("mb-6"),
				H2(
					Class("text-lg font-semibold text-gray-900 mb-4"),
					Text("AI-Generated Curriculum"),
				),
				Div(
					Class("bg-purple-50 border border-purple-200 rounded-lg p-6"),
					Pre(
						Class("whitespace-pre-wrap text-purple-800 leading-relaxed font-sans text-base"),
						Text(note.AiCurriculum),
					),
				),
			),
		),

		// AI Processing indicator
		If(note.AiProcessing,
			Div(
				Class("mb-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg"),
				Div(
					Class("flex items-center gap-3"),
					Div(
						Class("animate-spin rounded-full h-4 w-4 border-b-2 border-yellow-600"),
					),
					Text("AI is processing this note to generate a curriculum..."),
				),
			),
		),

		// JavaScript for share functionality
		Script(
			Raw(fmt.Sprintf(`
				function copyShareLink(token) {
					const shareUrl = '%s/share/' + token;
					navigator.clipboard.writeText(shareUrl).then(function() {
						alert('Share link copied to clipboard!');
					}, function(err) {
						console.error('Could not copy text: ', err);
					});
				}
			`, r.Config.App.Host)),
		),
	)

	return r.Render(layouts.Primary, content)
}

// EditNote displays the edit note form
func EditNote(ctx echo.Context, form *forms.EditNote) error {
	r := ui.NewRequest(ctx)
	r.Title = "Edit Note"

	content := Div(
		Class("max-w-4xl mx-auto"),

		// Form
		form.Render(r),
	)

	return r.Render(layouts.Primary, content)
}

// noteCard creates a card component for displaying a note in the list
func noteCard(r *ui.Request, note *ent.Note) Node {
	return Div(
		Class("bg-white rounded-lg border border-gray-200 hover:border-gray-300 hover:shadow-md transition-all duration-200"),
		A(
			Href(r.Path(routenames.Notes+".view", note.ID)),
			Class("block p-6 h-full"),

			// Header
			Div(
				Class("flex justify-between items-start mb-3"),
				H3(
					Class("text-lg font-semibold text-gray-900 line-clamp-2"),
					Text(note.Title),
				),
				Div(
					Class("flex-shrink-0 ml-2"),
					If(string(note.Visibility) == "public",
						Span(
							Class("inline-flex items-center px-2 py-1 rounded-full text-xs bg-green-100 text-green-800"),
							Text("Public"),
						),
					),
					If(string(note.Visibility) == "private",
						Span(
							Class("inline-flex items-center px-2 py-1 rounded-full text-xs bg-gray-100 text-gray-800"),
							Text("Private"),
						),
					),
				),
			),

			// Description
			If(note.Description != "",
				P(
					Class("text-gray-600 text-sm mb-4 line-clamp-3"),
					Text(note.Description),
				),
			),

			// Footer
			Div(
				Class("flex justify-between items-center text-xs text-gray-500 mt-auto"),
				Span(
					Text(note.UpdatedAt.Format("Jan 2, 2006")),
				),
				If(note.AiProcessing,
					Span(
						Class("inline-flex items-center text-yellow-600"),
						Div(
							Class("animate-spin rounded-full h-3 w-3 border-b border-yellow-600 mr-1"),
						),
						Text("AI Processing"),
					),
				),
				If(!note.AiProcessing && note.AiCurriculum != "",
					Span(
						Class("text-purple-600"),
						Text("AI Ready"),
					),
				),
			),
		),
	)
}