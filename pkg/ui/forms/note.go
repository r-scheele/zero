package forms

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/r-scheele/zero/pkg/form"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// formatFileSize converts bytes to a human-readable string
func formatFileSize(bytes int64) string {
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

// CreateNote represents the form for creating a new note
type CreateNote struct {
	Title           string   `form:"title" validate:"required,min=1,max=200"`
	Description     string   `form:"description" validate:"max=500"`
	Content         string   `form:"content"`
	Visibility      string   `form:"visibility" validate:"oneof=private public"`
	PermissionLevel string   `form:"permission_level" validate:"oneof=read_only read_write read_write_approval"`
	ResourceURLs    []string `form:"resource_urls"`
	MaxFileSize  int64
	MaxTotalSize int64
	MaxFiles     int
	form.Submission
}

// Render renders the create note form
func (f *CreateNote) Render(r *ui.Request) Node {
	return Form(
		ID("create-note"),
		Method(http.MethodPost),
		Attr("hx-post", r.Path(routenames.Notes+".create")),
		Class("space-y-6"),
		FlashMessages(r),

		// Title field
		InputField(InputFieldParams{
			Form:      f,
			FormField: "Title",
			Name:      "title",
			InputType: "text",
			Label:     "Title",
			Value:     f.Title,
			Help:      "Give your note a descriptive title",
		}),

		// Description field
		TextareaField(TextareaFieldParams{
			Form:      f,
			FormField: "Description",
			Name:      "description",
			Label:     "Description (Optional)",
			Value:     f.Description,
			Help:      "Brief description of what this note covers",
		}),

		// Content field
		TextareaField(TextareaFieldParams{
			Form:      f,
			FormField: "Content",
			Name:      "content",
			Label:     "Content",
			Value:     f.Content,
			Help:      "Main content of your note",
		}),

		// Visibility settings
		Div(
			Class("grid grid-cols-1 md:grid-cols-2 gap-4"),

			// Visibility field
			Div(
				Label(
					For("visibility"),
					Class("block text-sm font-medium text-gray-700 mb-2"),
					Text("Visibility"),
				),
				Select(
					ID("visibility"),
					Name("visibility"),
					Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"),
					Option(
						Value("private"),
						If(f.Visibility == "private" || f.Visibility == "", Selected()),
						Text("Private"),
					),
					Option(
						Value("public"),
						If(f.Visibility == "public", Selected()),
						Text("Public"),
					),
				),
				P(
					Class("text-sm text-gray-500 mt-1"),
					Text("Private notes are only visible to you"),
				),
			),

			// Permission level field (only shown for public notes)
			Div(
				ID("permission-level-container"),
				Class(func() string {
					if f.Visibility == "public" {
						return "block"
					}
					return "hidden"
				}()),
				Label(
					For("permission_level"),
					Class("block text-sm font-medium text-gray-700 mb-2"),
					Text("Permission Level"),
				),
				Select(
					ID("permission_level"),
					Name("permission_level"),
					Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"),
					Option(
						Value("read_only"),
						If(f.PermissionLevel == "read_only" || f.PermissionLevel == "", Selected()),
						Text("Read Only"),
					),
					Option(
						Value("read_write"),
						If(f.PermissionLevel == "read_write", Selected()),
						Text("Read & Write"),
					),
					Option(
						Value("read_write_approval"),
						If(f.PermissionLevel == "read_write_approval", Selected()),
						Text("Read & Write (Approval Required)"),
					),
				),
				P(
					Class("text-sm text-gray-500 mt-1"),
					Text("Controls what others can do with your public note"),
				),
			),
		),

		// Resource upload section
		Div(
			Class("space-y-4"),
			H3(
				Class("text-lg font-medium text-gray-900 mb-4"),
				Text("📎 Add Resources"),
			),
			P(
				Class("text-sm text-gray-600 mb-6"),
				Text("Upload files, add YouTube links, or attach documents to enhance your note with AI-powered curriculum generation."),
			),

			// File upload area - Blue theme
			Div(
				ID("file-upload-area"),
				Class("relative border-2 border-dashed border-blue-300 bg-blue-50 rounded-lg p-6 hover:border-blue-400 hover:bg-blue-100 transition-colors cursor-pointer"),
				Div(
					Class("flex items-center mb-3"),
					Div(
						Class("w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center mr-3"),
						Text("📁"),
					),
					H4(
						Class("text-md font-medium text-blue-900"),
						Text("Upload Files"),
					),
				),
				// Drag and drop visual cue
				Div(
					Class("text-center mb-4"),
					Div(
						Class("text-3xl mb-2"),
						Text("📤"),
					),
					P(
						Class("text-blue-700 font-medium mb-1"),
						Text("Drag and drop files here"),
					),
					P(
						Class("text-blue-600 text-sm"),
						Text("or click to browse files"),
					),
				),
				Input(
					Type("file"),
					ID("file-upload"),
					Name("files"),
					Multiple(),
					Accept(".pdf,.doc,.docx,.txt,.jpg,.jpeg,.png,.gif,.mp4,.avi,.mov"),
					Class("absolute inset-0 w-full h-full opacity-0 cursor-pointer"),
				),
				P(
					Class("text-xs text-blue-600 mt-2 text-center"),
					Text(fmt.Sprintf("Supported: PDF, DOC, images, videos (max %s each, %s total, %d files max)", 
					formatFileSize(f.MaxFileSize), 
					formatFileSize(f.MaxTotalSize), 
					f.MaxFiles)),
				),
			),

			// Multiple URL inputs for YouTube/web links - Green theme
			Div(
				Class("border-2 border-dashed border-green-300 bg-green-50 rounded-lg p-6"),
				Div(
					Class("flex items-center mb-3"),
					Div(
						Class("w-8 h-8 bg-green-100 rounded-full flex items-center justify-center mr-3"),
						Text("🔗"),
					),
					H4(
						Class("text-md font-medium text-green-900"),
						Text("Add Links"),
					),
				),
				// New URL input form
				Div(
					Class("flex gap-2 mb-3"),
					Input(
						Type("url"),
						ID("new-url-input"),
						Placeholder("https://youtube.com/watch?v=... or any URL"),
						Class("flex-1 px-3 py-2 border border-green-300 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500 bg-white"),
					),
					Button(
						Type("button"),
						ID("add-url-btn"),
						Class("px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors font-medium"),
						Text("Add Link"),
					),
				),
				
				// Display previously added URLs (read-only) - now below the input form
				Div(
					ID("added-urls-display"),
					Class("space-y-2 mb-4"),
					Group(func() []Node {
						var nodes []Node
						for i, url := range f.ResourceURLs {
							if strings.TrimSpace(url) != "" {
								nodes = append(nodes, Div(
									Class("flex items-center gap-3 p-3 bg-gray-50 border border-gray-200 rounded-md group hover:bg-gray-100 transition-colors"),
									Div(
										Class("flex-1 text-sm text-gray-700 break-all"),
										Text(url),
									),
									// Hidden input to preserve the URL value
									Input(
										Type("hidden"),
										Name("resource_urls"),
										Value(url),
									),
									Button(
											Type("button"),
											Class("px-3 py-2 bg-red-300 text-red-800 text-sm font-bold rounded-full hover:bg-red-400 hover:text-red-900 transition-all duration-200 shadow-md border-2 border-red-400 remove-added-url-btn hover:scale-110 hover:shadow-lg"),
											Attr("data-url-index", fmt.Sprintf("%d", i)),
											Attr("title", "Click to remove this link"),
											Text("×"),
										),
								))
							}
						}
						return nodes
					}()),
				),
				P(
					Class("text-xs text-green-600 mt-2"),
					Text("Add YouTube videos, articles, or any web resources"),
				),
			),
		),

		// Submit buttons
		ControlGroup(
			FormButton(ColorPrimary, "Create Note"),
			A(
				Href(r.Path(routenames.Notes)),
				Class("ml-3 px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 transition-colors"),
				Text("Cancel"),
			),
		),

		CSRF(r),

		// Inline JavaScript for note form functionality
		Script(
			Attr("type", "text/javascript"),
			Raw(fmt.Sprintf(`
				// Visibility handling for permission levels
				function updatePermissionVisibility() {
					const visibilitySelect = document.getElementById('visibility');
					const permissionContainer = document.getElementById('permission-level-container');
					
					if (visibilitySelect && permissionContainer) {
						if (visibilitySelect.value === 'public') {
							permissionContainer.classList.remove('hidden');
							permissionContainer.classList.add('block');
						} else {
							permissionContainer.classList.add('hidden');
							permissionContainer.classList.remove('block');
						}
					}
				}

				// File upload validation
				function initializeFileValidation(maxFiles, maxFileSize, maxTotalSize) {
					const fileInput = document.getElementById('file-upload');
					if (!fileInput) return;
					
					fileInput.addEventListener('change', function(e) {
						const files = e.target.files;
						
						if (files.length > maxFiles) {
							alert('Maximum ' + maxFiles + ' files allowed. Please select fewer files.');
							e.target.value = '';
							return;
						}
						
						let totalSize = 0;
						for (let i = 0; i < files.length; i++) {
							if (files[i].size > maxFileSize) {
								const maxSizeMB = Math.round(maxFileSize / (1024 * 1024));
								alert('File "' + files[i].name + '" exceeds ' + maxSizeMB + 'MB limit. Please choose a smaller file.');
								e.target.value = '';
								return;
							}
							totalSize += files[i].size;
						}
						
						if (totalSize > maxTotalSize) {
							const maxTotalMB = Math.round(maxTotalSize / (1024 * 1024));
							alert('Total file size exceeds ' + maxTotalMB + 'MB limit. Please select smaller files.');
							e.target.value = '';
							return;
						}
					});
				}

				// URL management functionality
				function addUrlInput() {
					const newUrlInput = document.getElementById('new-url-input');
					const displayContainer = document.getElementById('url-display-container');
					
					if (!newUrlInput || !displayContainer) return;
					
					const urlValue = newUrlInput.value.trim();
					if (!urlValue) return;
					
					const urlIndex = Date.now();
					const urlDisplay = document.createElement('div');
					urlDisplay.className = 'flex items-center gap-3 p-3 bg-gray-50 border border-gray-200 rounded-md group hover:bg-gray-100 transition-colors';
					
					const urlTextDiv = document.createElement('div');
					urlTextDiv.className = 'flex-1 text-sm text-gray-700 break-all';
					urlTextDiv.textContent = urlValue;
					
					const hiddenInput = document.createElement('input');
					hiddenInput.type = 'hidden';
					hiddenInput.name = 'resource_urls';
					hiddenInput.value = urlValue;
					
					const removeButton = document.createElement('button');
					removeButton.type = 'button';
					removeButton.className = 'px-3 py-2 bg-red-300 text-red-800 text-sm font-bold rounded-full hover:bg-red-400 hover:text-red-900 transition-all duration-200 shadow-md border-2 border-red-400 remove-added-url-btn hover:scale-110 hover:shadow-lg';
					removeButton.setAttribute('title', 'Click to remove this link');
					removeButton.setAttribute('data-url-index', urlIndex);
					removeButton.textContent = '×';
					
					urlDisplay.appendChild(urlTextDiv);
					urlDisplay.appendChild(hiddenInput);
					urlDisplay.appendChild(removeButton);
					
					removeButton.addEventListener('click', function() {
						urlDisplay.remove();
					});
					
					displayContainer.appendChild(urlDisplay);
					newUrlInput.value = '';
					newUrlInput.focus();
				}

				// Initialize form functionality
				document.addEventListener('DOMContentLoaded', function() {
					// Initialize visibility handling
					updatePermissionVisibility();
					
					// Handle visibility changes
					const visibilitySelect = document.getElementById('visibility');
					if (visibilitySelect) {
						visibilitySelect.addEventListener('change', function() {
							updatePermissionVisibility();
						});
					}
					
					// Initialize file validation
					initializeFileValidation(%d, %d, %d);
					
					// Add event listener to the "Add Link" button
					const addBtn = document.getElementById('add-url-btn');
					if (addBtn) {
						addBtn.addEventListener('click', function(e) {
							e.preventDefault();
							addUrlInput();
						});
					}
					
					// Add event listener for Enter key on the input field
					const newUrlInput = document.getElementById('new-url-input');
					if (newUrlInput) {
						newUrlInput.addEventListener('keypress', function(e) {
							if (e.key === 'Enter') {
								e.preventDefault();
								addUrlInput();
							}
						});
					}
					
					// Add event listeners to existing remove buttons
					const removeButtons = document.querySelectorAll('.remove-added-url-btn');
					removeButtons.forEach(function(btn) {
						btn.addEventListener('click', function() {
							btn.closest('div').remove();
						});
					});
				});
			`, f.MaxFiles, f.MaxFileSize, f.MaxTotalSize)),
		),
	)
}

// EditNote represents the form for editing an existing note
type EditNote struct {
	ID              int      `form:"id"`
	Title           string   `form:"title" validate:"required,min=1,max=200"`
	Description     string   `form:"description" validate:"max=500"`
	Content         string   `form:"content"`
	Visibility      string   `form:"visibility" validate:"oneof=private public"`
	PermissionLevel string   `form:"permission_level" validate:"oneof=read_only read_write read_write_approval"`
	ResourceURLs    []string `form:"resource_urls"`
	MaxFileSize     int64
	MaxTotalSize    int64
	MaxFiles        int
	form.Submission
}

// Render renders the edit note form
func (f *EditNote) Render(r *ui.Request) Node {
	return Form(
		ID("edit-note"),
		Method(http.MethodPost),
		Attr("hx-post", r.Path(routenames.Notes+".edit", f.ID)),
		Class("space-y-6"),
		FlashMessages(r),

		// Hidden ID field
		Input(
			Type("hidden"),
			Name("id"),
			Value(fmt.Sprintf("%d", f.ID)),
		),

		// Title field
		InputField(InputFieldParams{
			Form:      f,
			FormField: "Title",
			Name:      "title",
			InputType: "text",
			Label:     "Title",
			Value:     f.Title,
		}),

		// Description field
		TextareaField(TextareaFieldParams{
			Form:      f,
			FormField: "Description",
			Name:      "description",
			Label:     "Description (Optional)",
			Value:     f.Description,
		}),

		// Content field
		TextareaField(TextareaFieldParams{
			Form:      f,
			FormField: "Content",
			Name:      "content",
			Label:     "Content",
			Value:     f.Content,
		}),

		// Resource upload section
		Div(
			Class("space-y-4"),
			H3(
				Class("text-lg font-medium text-gray-900 mb-4"),
				Text("📎 Edit Resources"),
			),
			P(
				Class("text-sm text-gray-600 mb-6"),
				Text("Upload files, add YouTube links, or attach documents to enhance your note."),
			),

			// File upload area
			Div(
				ID("file-upload-area"),
				Class("relative border-2 border-dashed border-blue-300 bg-blue-50 rounded-lg p-6 hover:border-blue-400 hover:bg-blue-100 transition-colors cursor-pointer"),
				Div(
					Class("flex items-center mb-3"),
					Div(
						Class("w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center mr-3"),
						Text("📁"),
					),
					H4(
						Class("text-md font-medium text-blue-900"),
						Text("Upload Files"),
					),
				),
				Div(
					Class("text-center mb-4"),
					Div(
						Class("text-3xl mb-2"),
						Text("📤"),
					),
					P(
						Class("text-blue-700 font-medium mb-1"),
						Text("Drag and drop files here"),
					),
					P(
						Class("text-blue-600 text-sm"),
						Text("or click to browse files"),
					),
				),
				Input(
					Type("file"),
					ID("file-upload"),
					Name("files"),
					Multiple(),
					Accept(".pdf,.doc,.docx,.txt,.jpg,.jpeg,.png,.gif,.mp4,.avi,.mov"),
					Class("absolute inset-0 w-full h-full opacity-0 cursor-pointer"),
				),
				P(
					Class("text-xs text-blue-600 mt-2 text-center"),
					Text(fmt.Sprintf("Supported: PDF, DOC, images, videos (max %s each, %s total, %d files max)", 
						formatFileSize(f.MaxFileSize), formatFileSize(f.MaxTotalSize), f.MaxFiles)),
				),
			),

			// URL input section
			Div(
				Class("bg-green-50 border border-green-200 rounded-lg p-4"),
				Div(
					Class("flex items-center mb-3"),
					Div(
						Class("w-8 h-8 bg-green-100 rounded-full flex items-center justify-center mr-3"),
						Text("🔗"),
					),
					H4(
						Class("text-md font-medium text-green-900"),
						Text("Add Web Links"),
					),
				),
				Div(
					Class("flex gap-2"),
					Input(
						Type("url"),
						ID("new-url-input"),
						Placeholder("https://example.com or YouTube link"),
						Class("flex-1 px-3 py-2 border border-green-300 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500"),
					),
					Button(
						Type("button"),
						ID("add-url-btn"),
						Class("px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"),
						Text("Add"),
					),
				),
			),

			// Display added URLs
			Div(
				ID("added-urls-container"),
				Class("space-y-2"),
				Group(func() []Node {
					var nodes []Node
					for i, url := range f.ResourceURLs {
						if strings.TrimSpace(url) != "" {
							nodes = append(nodes, Div(
								Class("flex items-center justify-between bg-gray-50 border border-gray-200 rounded-md p-3"),
								Span(
									Class("flex-1 text-sm text-gray-700 break-all"),
									Text(url),
								),
								Input(
									Type("hidden"),
									Name("resource_urls"),
									Value(url),
								),
								Button(
									Type("button"),
									Class("px-3 py-2 bg-red-300 text-red-800 text-sm font-bold rounded-full hover:bg-red-400 hover:text-red-900 transition-all duration-200 shadow-md border-2 border-red-400 remove-added-url-btn hover:scale-110 hover:shadow-lg"),
									Attr("data-url-index", fmt.Sprintf("%d", i)),
									Attr("title", "Click to remove this link"),
									Text("×"),
								),
							))
						}
					}
					return nodes
				}()),
			),
			P(
				Class("text-xs text-green-600 mt-2"),
				Text("Add YouTube videos, articles, or any web resources"),
			),
		),

		// Visibility and permission settings
		Div(
			Class("grid grid-cols-1 md:grid-cols-2 gap-4"),

			// Visibility field
			Div(
				Label(
					For("visibility"),
					Class("block text-sm font-medium text-gray-700 mb-2"),
					Text("Visibility"),
				),
				Select(
					ID("visibility"),
					Name("visibility"),
					Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"),
					Option(
						Value("private"),
						If(f.Visibility == "private", Selected()),
						Text("Private"),
					),
					Option(
						Value("public"),
						If(f.Visibility == "public", Selected()),
						Text("Public"),
					),
				),
			),

			// Permission level field (only shown for public notes)
			Div(
				ID("permission-level-container"),
				Class(func() string {
					if f.Visibility == "public" {
						return "block"
					}
					return "hidden"
				}()),
				Label(
					For("permission_level"),
					Class("block text-sm font-medium text-gray-700 mb-2"),
					Text("Permission Level"),
				),
				Select(
					ID("permission_level"),
					Name("permission_level"),
					Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"),
					Option(
						Value("read_only"),
						If(f.PermissionLevel == "read_only", Selected()),
						Text("Read Only"),
					),
					Option(
						Value("read_write"),
						If(f.PermissionLevel == "read_write", Selected()),
						Text("Read & Write"),
					),
					Option(
						Value("read_write_approval"),
						If(f.PermissionLevel == "read_write_approval", Selected()),
						Text("Read & Write (Approval Required)"),
					),
				),
			),
		),

		// Submit buttons
		ControlGroup(
			FormButton(ColorPrimary, "Update Note"),
			A(
				Href(r.Path(routenames.Notes+".view", f.ID)),
				Class("ml-3 px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 transition-colors"),
				Text("Cancel"),
			),
		),

		CSRF(r),

		// Inline JavaScript for note form functionality
		Script(
			Attr("type", "text/javascript"),
			Raw(fmt.Sprintf(`
				// Visibility handling for permission levels
				function updatePermissionVisibility() {
					const visibilitySelect = document.getElementById('visibility');
					const permissionContainer = document.getElementById('permission-level-container');
					
					if (visibilitySelect && permissionContainer) {
						if (visibilitySelect.value === 'public') {
							permissionContainer.classList.remove('hidden');
							permissionContainer.classList.add('block');
						} else {
							permissionContainer.classList.add('hidden');
							permissionContainer.classList.remove('block');
						}
					}
				}

				// File upload validation
				function initializeFileValidation(maxFiles, maxFileSize, maxTotalSize) {
					const fileInput = document.getElementById('file-upload');
					if (!fileInput) return;
					
					fileInput.addEventListener('change', function(e) {
						const files = e.target.files;
						
						if (files.length > maxFiles) {
							alert('Maximum ' + maxFiles + ' files allowed. Please select fewer files.');
							e.target.value = '';
							return;
						}
						
						let totalSize = 0;
						for (let i = 0; i < files.length; i++) {
							if (files[i].size > maxFileSize) {
								const maxSizeMB = Math.round(maxFileSize / (1024 * 1024));
								alert('File "' + files[i].name + '" exceeds ' + maxSizeMB + 'MB limit. Please choose a smaller file.');
								e.target.value = '';
								return;
							}
							totalSize += files[i].size;
						}
						
						if (totalSize > maxTotalSize) {
							const maxTotalMB = Math.round(maxTotalSize / (1024 * 1024));
							alert('Total file size exceeds ' + maxTotalMB + 'MB limit. Please select smaller files.');
							e.target.value = '';
							return;
						}
					});
				}

				// URL management functionality
				function addUrlInput() {
					const newUrlInput = document.getElementById('new-url-input');
					const displayContainer = document.getElementById('url-display-container');
					
					if (!newUrlInput || !displayContainer) return;
					
					const urlValue = newUrlInput.value.trim();
					if (!urlValue) return;
					
					const urlIndex = Date.now();
					const urlDisplay = document.createElement('div');
					urlDisplay.className = 'flex items-center gap-3 p-3 bg-gray-50 border border-gray-200 rounded-md group hover:bg-gray-100 transition-colors';
					
					const urlTextDiv = document.createElement('div');
					urlTextDiv.className = 'flex-1 text-sm text-gray-700 break-all';
					urlTextDiv.textContent = urlValue;
					
					const hiddenInput = document.createElement('input');
					hiddenInput.type = 'hidden';
					hiddenInput.name = 'resource_urls';
					hiddenInput.value = urlValue;
					
					const removeButton = document.createElement('button');
					removeButton.type = 'button';
					removeButton.className = 'px-3 py-2 bg-red-300 text-red-800 text-sm font-bold rounded-full hover:bg-red-400 hover:text-red-900 transition-all duration-200 shadow-md border-2 border-red-400 remove-added-url-btn hover:scale-110 hover:shadow-lg';
					removeButton.setAttribute('title', 'Click to remove this link');
					removeButton.setAttribute('data-url-index', urlIndex);
					removeButton.textContent = '×';
					
					urlDisplay.appendChild(urlTextDiv);
					urlDisplay.appendChild(hiddenInput);
					urlDisplay.appendChild(removeButton);
					
					removeButton.addEventListener('click', function() {
						urlDisplay.remove();
					});
					
					displayContainer.appendChild(urlDisplay);
					newUrlInput.value = '';
					newUrlInput.focus();
				}

				// Initialize form functionality
				document.addEventListener('DOMContentLoaded', function() {
					// Initialize visibility handling
					updatePermissionVisibility();
					
					// Handle visibility changes
					const visibilitySelect = document.getElementById('visibility');
					if (visibilitySelect) {
						visibilitySelect.addEventListener('change', function() {
							updatePermissionVisibility();
						});
					}
					
					// Initialize file validation
					initializeFileValidation(%d, %d, %d);
					
					// Add event listener to the "Add Link" button
					const addBtn = document.getElementById('add-url-btn');
					if (addBtn) {
						addBtn.addEventListener('click', function(e) {
							e.preventDefault();
							addUrlInput();
						});
					}
					
					// Add event listener for Enter key on the input field
					const newUrlInput = document.getElementById('new-url-input');
					if (newUrlInput) {
						newUrlInput.addEventListener('keypress', function(e) {
							if (e.key === 'Enter') {
								e.preventDefault();
								addUrlInput();
							}
						});
					}
					
					// Add event listeners to existing remove buttons
					const removeButtons = document.querySelectorAll('.remove-added-url-btn');
					removeButtons.forEach(function(btn) {
						btn.addEventListener('click', function() {
							btn.closest('div').remove();
						});
					});
				});
			`, f.MaxFiles, f.MaxFileSize, f.MaxTotalSize)),
		),
	)
}
