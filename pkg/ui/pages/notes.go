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
		Class("max-w-4xl mx-auto space-y-8"),

		// Header with actions
		Div(
			Class("flex justify-between items-start"),
			Div(
				Class("flex-1"),
				H1(
					Class("text-3xl font-bold text-gray-900 mb-2"),
					Text(note.Title),
				),
				Div(
					Class("flex items-center space-x-4 text-sm text-gray-600"),
					Span(
						Text("By "),
						If(note.Edges.Owner != nil,
							Strong(Text(note.Edges.Owner.Name)),
						),
					),
					Span(
						Text("•"),
					),
					Span(
						Text(note.CreatedAt.Format("Jan 2, 2006")),
					),
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

			// Actions (only for owner)
			If(isOwner,
				Div(
					Class("flex items-center space-x-3"),

					// Share button
					Button(
						Type("button"),
						Class("inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm text-gray-700 bg-white hover:bg-gray-50"),
						Attr("onclick", fmt.Sprintf("copyShareLink('%s')", note.ShareToken)),
						Text("Share"),
					),

					// Edit button
					A(
						Href(r.Path(routenames.Notes+".edit", note.ID)),
						Class("inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm text-gray-700 bg-white hover:bg-gray-50"),
						Text("Edit"),
					),

					// Delete button
					Form(
						Method("POST"),
						Action(r.Path(routenames.Notes+".delete", note.ID)),
						Attr("onsubmit", "return confirm('Are you sure you want to delete this note?')"),
						CSRF(r),
						Button(
							Type("submit"),
							Class("inline-flex items-center px-3 py-2 border border-red-300 rounded-md text-sm text-red-700 bg-white hover:bg-red-50"),
							Text("Delete"),
						),
					),
				),
			),
		),

		// Description
		If(note.Description != "",
			Div(
				Class("bg-blue-50 border-l-4 border-blue-400 p-4 rounded-r-lg"),
				P(
					Class("text-blue-800"),
					Text(note.Description),
				),
			),
		),

		// Content
		If(note.Content != "",
			Div(
				Class("bg-white rounded-lg border border-gray-200 p-6"),
				H2(
					Class("text-xl font-semibold text-gray-900 mb-4"),
					Text("Content"),
				),
				Div(
					Class("prose max-w-none"),
					Pre(
						Class("whitespace-pre-wrap text-gray-700 leading-relaxed"),
						Text(note.Content),
					),
				),
			),
		),

		// AI Curriculum (if available)
		If(note.AiCurriculum != "",
			Div(
				Class("bg-gradient-to-r from-purple-50 to-pink-50 rounded-lg border border-purple-200 p-6"),
				H2(
					Class("text-xl font-semibold text-purple-900 mb-4 flex items-center"),
					Text("AI-Generated Curriculum"),
				),
				Div(
					Class("prose max-w-none"),
					Pre(
						Class("whitespace-pre-wrap text-purple-800 leading-relaxed"),
						Text(note.AiCurriculum),
					),
				),
			),
		),

		// AI Processing indicator
		If(note.AiProcessing,
			Div(
				Class("bg-yellow-50 border border-yellow-200 rounded-lg p-4"),
				Div(
					Class("flex items-center"),
					Div(
						Class("animate-spin rounded-full h-4 w-4 border-b-2 border-yellow-600 mr-3"),
					),
					Text("AI is processing this note to generate a curriculum..."),
				),
			),
		),

		// JavaScript for share functionality
		Script(
			Text(fmt.Sprintf(`
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

		// Header
		Div(
			Class("mb-8"),
			H1(
				Class("text-3xl font-bold text-gray-900 mb-2"),
				Text("Edit Note"),
			),
		),

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