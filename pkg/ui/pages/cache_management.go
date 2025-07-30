package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func CacheManagement(ctx echo.Context, f *forms.Cache) error {
	r := ui.NewRequest(ctx)
	r.Title = "Cache Management - Admin"

	content := cacheManagementContent(f)

	return r.Render(layouts.Admin, content)
}

func cacheManagementContent(f *forms.Cache) Node {
	return Div(
		Class("space-y-6"),
		
		// Header
		Div(
			Class("border-b border-gray-200 pb-4"),
			H1(
				Class("text-2xl font-bold text-gray-900"),
				Text("Cache Management"),
			),
			P(
				Class("mt-2 text-sm text-gray-600"),
				Text("Manage application cache to optimize performance and clear stale data."),
			),
		),
		
		// Cache Statistics
		Div(
			Class("grid grid-cols-1 md:grid-cols-3 gap-6"),
			
			// User Cache
			Div(
				Class("bg-white overflow-hidden shadow rounded-lg"),
				Div(
					Class("p-5"),
					Div(
						Class("flex items-center"),
						Div(
							Class("flex-shrink-0"),
							Div(
								Class("w-8 h-8 bg-blue-500 rounded-md flex items-center justify-center text-white"),
								Text("👤"),
							),
						),
						Div(
							Class("ml-5 w-0 flex-1"),
							Dl(
								Dt(
									Class("text-sm font-medium text-gray-500 truncate"),
									Text("User Sessions"),
								),
								Dd(
									Class("text-lg font-medium text-gray-900"),
									Text("Active"),
								),
							),
						),
					),
				),
			),
			
			// Note Cache
			Div(
				Class("bg-white overflow-hidden shadow rounded-lg"),
				Div(
					Class("p-5"),
					Div(
						Class("flex items-center"),
						Div(
							Class("flex-shrink-0"),
							Div(
								Class("w-8 h-8 bg-green-500 rounded-md flex items-center justify-center text-white"),
								Text("📝"),
							),
						),
						Div(
							Class("ml-5 w-0 flex-1"),
							Dl(
								Dt(
									Class("text-sm font-medium text-gray-500 truncate"),
									Text("Note Data"),
								),
								Dd(
									Class("text-lg font-medium text-gray-900"),
									Text("Cached"),
								),
							),
						),
					),
				),
			),
			
			// Response Cache
			Div(
				Class("bg-white overflow-hidden shadow rounded-lg"),
				Div(
					Class("p-5"),
					Div(
						Class("flex items-center"),
						Div(
							Class("flex-shrink-0"),
							Div(
								Class("w-8 h-8 bg-purple-500 rounded-md flex items-center justify-center text-white"),
								Text("🌐"),
							),
						),
						Div(
							Class("ml-5 w-0 flex-1"),
							Dl(
								Dt(
									Class("text-sm font-medium text-gray-500 truncate"),
									Text("HTTP Responses"),
								),
								Dd(
									Class("text-lg font-medium text-gray-900"),
									Text("Optimized"),
								),
							),
						),
					),
				),
			),
		),
		
		// Cache Actions
		Div(
			Class("bg-white shadow rounded-lg"),
			Div(
				Class("px-4 py-5 sm:p-6"),
				H3(
					Class("text-lg leading-6 font-medium text-gray-900 mb-4"),
					Text("Cache Actions"),
				),
				
				// Clear All Cache
				Div(
					Class("mb-6"),
					Div(
						Class("flex items-center justify-between p-4 border border-gray-200 rounded-lg"),
						Div(
							H4(
								Class("text-sm font-medium text-gray-900"),
								Text("Clear All Cache"),
							),
							P(
								Class("text-sm text-gray-500 mt-1"),
								Text("Remove all cached data including user sessions, note data, and HTTP responses."),
							),
						),
						Button(
							Type("button"),
							Class("inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"),
							Text("Clear All"),
							Attr("onclick", "clearAllCache()"),
						),
					),
				),
				
				// Clear Specific Pattern
				Div(
					Class("border border-gray-200 rounded-lg p-4"),
					H4(
						Class("text-sm font-medium text-gray-900 mb-3"),
						Text("Clear Specific Cache Pattern"),
					),
					Form(
						Class("flex gap-3"),
						Attr("onsubmit", "clearCachePattern(event)"),
						Div(
							Class("flex-1"),
							Input(
								Type("text"),
								Name("pattern"),
								ID("cache-pattern"),
								Class("block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"),
								Placeholder("e.g., user:, note_likes_count:, response_cache:"),
							),
						),
						Button(
							Type("submit"),
							Class("inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"),
							Text("Clear Pattern"),
						),
					),
				),
			),
		),
		
		// JavaScript for cache management
		Script(
			Text(`
				function clearAllCache() {
					if (confirm('Are you sure you want to clear all cache? This action cannot be undone.')) {
						fetch('/admin/cache/clear', {
							method: 'POST',
							headers: {
								'Content-Type': 'application/json',
							},
						})
						.then(response => response.json())
						.then(data => {
							if (data.success) {
								alert('Cache cleared successfully!');
								location.reload();
							} else {
								alert('Error clearing cache: ' + data.message);
							}
						})
						.catch(error => {
							console.error('Error:', error);
							alert('Error clearing cache');
						});
					}
				}
				
				function clearCachePattern(event) {
					event.preventDefault();
					const pattern = document.getElementById('cache-pattern').value;
					if (!pattern) {
						alert('Please enter a cache pattern');
						return;
					}
					
					const formData = new FormData();
					formData.append('pattern', pattern);
					
					fetch('/admin/cache/clear-pattern', {
						method: 'POST',
						body: formData,
					})
					.then(response => response.json())
					.then(data => {
						if (data.success) {
							alert('Cache pattern cleared successfully!');
							document.getElementById('cache-pattern').value = '';
						} else {
							alert('Error clearing cache pattern: ' + data.message);
						}
					})
					.catch(error => {
						console.error('Error:', error);
						alert('Error clearing cache pattern');
					});
				}
			`),
		),
	)
}