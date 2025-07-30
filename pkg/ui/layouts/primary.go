package layouts

import (
	"strings"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/cache"
	. "github.com/r-scheele/zero/pkg/ui/components"
	"github.com/r-scheele/zero/pkg/ui/icons"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Primary(r *ui.Request, content Node) Node {
	// Basic safety check
	if r == nil {
		return Doctype(
			HTML(
				Lang("en"),
				Body(
					Class("min-h-screen bg-gray-100"),
					Div(
						Class("flex items-center justify-center min-h-screen"),
						Text("System Error"),
					),
				),
			),
		)
	}

	return Doctype(
		HTML(
			Lang("en"),
			Data("theme", "light"),
			Head(
				Metatags(r),
				CSS(),
				Style(`
					.more-active {
						background-color: rgb(219 234 254) !important; /* bg-blue-100 */
					}
					.more-active .text-slate-500 {
						color: rgb(37 99 235) !important; /* text-blue-600 */
					}
					.more-active .group-hover\\:bg-slate-100 {
						background-color: rgb(219 234 254) !important; /* bg-blue-100 */
					}
				`),
				JS(),
			),
			Body(
				Class("min-h-screen bg-gradient-to-br from-slate-50 to-gray-100"),
				Attr("onclick", `
					if (!event.target.closest('#mobile_nav_modal') && !event.target.closest('#more_nav_button')) {
						const modal = document.getElementById('mobile_nav_modal');
						const button = document.getElementById('more_nav_button');
						if (modal && button) {
							modal.classList.add('hidden');
							button.classList.remove('more-active');
						}
					}
				`),
				// Use the same header for all users
				unifiedHeader(r),
				If(r.IsAuth,
				// Authenticated user layout with sidebar
				Group{
					userSidebar(r), // Fixed positioned sidebar
					Main(
					Class("bg-slate-50 min-h-screen lg:ml-80"), // Only apply left margin on large screens where sidebar is visible
						Div(
							Class("px-4 sm:px-6 lg:px-8 xl:px-12 py-6 lg:py-8 pb-20 lg:pb-8"), // Extra bottom padding on mobile for nav
							Div(
								Class("max-w-7xl mx-auto"),
								ID("main-content"),
								FlashMessages(r),
								Div(
									Class("space-y-6"),
									content,
								),
							),
						),
					),
				},
				),
				If(!r.IsAuth,
					// Non-authenticated user layout without sidebar
					Div(
						Class("min-h-screen"),
						Main(
							Class("w-full bg-slate-50"),
							Div(
								Class("px-4 sm:px-6 lg:px-8 xl:px-12 py-6 lg:py-8 pb-20 lg:pb-8"),
								Div(
									Class("max-w-7xl mx-auto"),
									ID("main-content"),
									FlashMessages(r),
									Div(
										Class("space-y-6"),
										content,
									),
								),
							),
						),
					),
				),
				// Use the same footer for all users
				unifiedFooter(r),
				If(r.IsAuth, userBottomNavigation(r)), // Mobile navigation only for auth users
				searchModal(r),
				If(r.IsAuth, mobileNavModal(r)), // Additional nav items for mobile
				HtmxListeners(r),
			),
		),
	)
}

func userHeader(r *ui.Request) Node {
	// Comprehensive safety check to prevent panics with corrupted request context
	if r == nil {
		return Div() // Return empty div if request is nil
	}

	// Additional safety check - if context is corrupted, return minimal header
	defer func() {
		if recover() != nil {
			// If panic occurs during rendering, return minimal header
		}
	}()

	return Header(
		Class("fixed top-0 left-0 right-0 z-50 bg-white/95 backdrop-blur-lg border-b border-slate-200/50 shadow-elegant-sm"),
		Div(
			Class("container mx-auto max-w-none lg:max-w-6xl xl:max-w-7xl px-4 sm:px-6 lg:px-8 xl:px-12 py-3 lg:py-4"),
			Div(
				Class("flex items-center justify-between"),
				// Left side - Logo and App title
				Div(
					Class("flex items-center gap-4"),
					A(
						Href("/"),
						Class("flex items-center gap-3 text-slate-900 hover:text-blue-600 transition-colors duration-300"),
						Div(
							Class("flex items-center justify-center w-10 h-10 bg-gradient-to-r from-blue-600 to-blue-700 rounded-xl shadow-lg"),
							Span(
								Class("text-white font-bold text-xl"),
								Text("Z"),
							),
						),
						Div(
							H1(
								Class("text-xl lg:text-2xl font-bold"),
								Text("Zero"),
							),
						),
					),
				),
				// Right side - User info and actions
				Div(
					Class("flex items-center gap-4"),
					// User info - only show if user is authenticated
					If(r.IsAuth && r.AuthUser != nil,
						Div(
							Class("hidden sm:flex items-center gap-2 text-sm text-slate-600"),
							icons.UserCircle(),
							Span(Text(func() string {
								if r.AuthUser.Name != "" {
									return r.AuthUser.Name
								}
								return "User"
							}())),
							If(r.IsAdmin, Span(Class("badge badge-sm bg-blue-100 text-blue-800 border-blue-200"), Text("Admin"))),
						),
					),
					// Logout button - only show if user is authenticated
					If(r.IsAuth && r.AuthUser != nil,
						A(
							Href(r.Path(routenames.Logout)),
							Class("btn btn-outline btn-sm gap-2"),
							icons.Exit(),
							Span(Class("hidden sm:inline"), Text("Logout")),
						),
					),
				),
			),
		),
	)
}

func userSidebar(r *ui.Request) Node {
	// Safety check to prevent panics with corrupted request context
	if r == nil || r.Context == nil {
		return Div() // Return empty div if request is corrupted
	}

	userMenuItem := func(icon Node, title, href string) Node {
		// Path matching with special handling for profile and notes sections
		isActive := r.CurrentPath == href
		
		// Special case for Profile section - highlight when in any profile-related page
		if href == "/profile" && (r.CurrentPath == "/profile" || r.CurrentPath == "/profile/edit" || 
			r.CurrentPath == "/profile/update" || r.CurrentPath == "/profile/picture" || 
			r.CurrentPath == "/profile/change-password" || r.CurrentPath == "/profile/deactivate") {
			isActive = true
		}
		
		// Special case for Notes section - highlight when in any notes-related page
		if href == "/notes" && (strings.HasPrefix(r.CurrentPath, "/notes")) {
			isActive = true
		}

		var linkClasses string
		var iconClasses string
		var textClasses string

		// Special styling for Sign Out button
		if title == "Sign Out" {
			linkClasses = "user-nav-item flex items-center gap-3 px-4 py-3 rounded-xl text-red-600 hover:bg-red-50 transition-all duration-200"
			iconClasses = "w-5 h-5 text-red-600"
			textClasses = "font-medium text-sm text-red-600"
		} else if isActive {
			linkClasses = "user-nav-item user-nav-active flex items-center gap-3 px-4 py-3 rounded-xl bg-blue-100 text-blue-700 border border-blue-200 transition-all duration-200"
			iconClasses = "w-5 h-5 text-blue-600"
			textClasses = "font-semibold text-sm text-blue-700"
		} else {
			linkClasses = "user-nav-item flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-slate-900 transition-all duration-200"
			iconClasses = "w-5 h-5 text-slate-500 group-hover:text-blue-600"
			textClasses = "font-medium text-sm group-hover:text-slate-900"
		}

		return Li(
			Style("padding: 0 !important; margin: 0 !important; list-style: none !important; font-size: inherit !important;"),
			A(
				Href(href),
				Class(linkClasses+" group"),
				Style("display: flex !important; align-items: center !important; gap: 0.75rem !important; padding: 0.75rem 1rem !important; border-radius: 0.75rem !important; text-decoration: none !important; font-size: 0.875rem !important; line-height: 1.25rem !important; transition: all 0.2s !important;"),
				Div(
					Class(iconClasses),
					Style("width: 1.25rem !important; height: 1.25rem !important; flex-shrink: 0 !important;"),
					icon,
				),
				Span(
					Class(textClasses),
					Style("font-size: 0.875rem !important; line-height: 1.25rem !important; font-weight: 500 !important;"),
					Text(title),
				),
			),
		)
	}

	header := func(text string) Node {
		return Li(
			Class("px-4 py-2 text-xs font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100 mb-2 mt-6 first:mt-0"),
			Text(text),
		)
	}

	return Aside(
		Class("hidden lg:block w-80 bg-white border-r border-slate-200 shadow-sm flex-shrink-0 overflow-y-auto user-sidebar"),
		Style("font-size: 14px !important; line-height: 1.5 !important; font-family: system-ui, -apple-system, sans-serif !important; position: fixed; top: 0; left: 0; height: 100vh; z-index: 40;"),
		Div(
			Class("p-6 h-full"),
			Style("padding-top: 5rem !important;"), // Add top padding to account for fixed header
			Div(
				Class("flex flex-col h-full"), // Full height container
				Style("min-height: calc(100vh - 200px);"), // Ensure minimum height for proper flex behavior
				Nav(
					Class("flex-1 flex flex-col"),
					Style("all: revert !important; font-family: inherit !important;"),
					Ul(
						Class("space-y-2 flex flex-col h-full flex-1"),
						Style("list-style: none !important; padding: 0 !important; margin: 0 !important; font-size: inherit !important;"),
						HxBoost(),
						// Main navigation items - show only Home if user is not verified
						If(r.AuthUser != nil && !r.AuthUser.Verified,
							// Unverified user - only show Home
							Div(
								Class("space-y-2"),
								userMenuItem(icons.Home(), "Home", "/home"),
							),
						),
						If(r.AuthUser != nil && r.AuthUser.Verified,
							// Verified user - show all navigation items
							Div(
								Class("space-y-2"),
								userMenuItem(icons.Home(), "Home", "/home"),
								userMenuItem(icons.CircleStack(), "Dashboard", "/dashboard"),
								userMenuItem(icons.Star(), "Quizzes", "/quizzes"),
								userMenuItem(icons.Archive(), "Notes", "/notes"),
								userMenuItem(icons.UserCircle(), "Profile", "/profile"),
								userMenuItem(icons.Document(), "Files", "/files"), // Use direct path
							),
						),
						// Visual separator line after regular nav items - only for verified users
						If(r.AuthUser != nil && r.AuthUser.Verified,
							Div(
								Class("my-4 border-t border-slate-200"),
							),
						),
						// Admin section (if admin and verified)
						If(r.IsAdmin && r.AuthUser != nil && r.AuthUser.Verified,
							Div(
								Class("space-y-2"),
								header("⚙️ Admin Tools"),
								userMenuItem(icons.Archive(), "Cache Management", "/cache"), // Use direct path
								userMenuItem(icons.CircleStack(), "Background Tasks", "/admin/tasks"),
								userMenuItem(icons.UserCircle(), "User Management", "/admin/entity/user"),
							),
						),
						// Spacer to push Sign Out to the very bottom of viewport
						Div(
							Class("flex-grow"),
							Style("min-height: calc(100vh - 400px); min-height: calc(100svh - 400px);"), // Maximize space to push button to bottom
						),
						// Sign Out at the absolute bottom with visual separation - always show for authenticated users
			Div(
				Class("mt-auto pt-6 border-t border-slate-200"),
				Style("position: absolute; bottom: 0; left: 0; right: 0; background: white; padding: 1.5rem; border-top: 1px solid rgb(226 232 240);"), // Absolute positioning at bottom
				A(
					Href(r.Path(routenames.Logout)),
					Class("user-nav-item flex items-center gap-3 px-4 py-3 rounded-xl text-red-600 hover:bg-red-50 transition-all duration-200"),
					Attr("hx-boost", "false"), // Disable HTMX boost for reliable logout
					Div(
						Class("w-5 h-5 text-red-600"),
						icons.Exit(),
					),
					Span(
						Class("font-medium text-sm text-red-600"),
						Text("Sign Out"),
					),
				),
			),
					),
				),
			),
		),
	)
}

func userBottomNavigation(r *ui.Request) Node {
	// Safety check to prevent panics with corrupted request context
	if r == nil || r.Context == nil {
		return Div() // Return empty div if request is corrupted
	}

	// Bottom navigation item
	userNavItem := func(icon Node, title, href string) Node {
		// Path matching with special handling for profile and notes sections
		isActive := r.CurrentPath == href
		
		// Special case for Profile section - highlight when in any profile-related page
		if href == "/profile" && (r.CurrentPath == "/profile" || r.CurrentPath == "/profile/edit" || 
			r.CurrentPath == "/profile/update" || r.CurrentPath == "/profile/picture" || 
			r.CurrentPath == "/profile/change-password" || r.CurrentPath == "/profile/deactivate") {
			isActive = true
		}
		
		// Special case for Notes section - highlight when in any notes-related page
		if href == "/notes" && (strings.HasPrefix(r.CurrentPath, "/notes")) {
			isActive = true
		}

		var iconContainerClass string
		var textClass string

		// Special styling for Sign Out button
		if title == "Sign Out" {
			if isActive {
				iconContainerClass = "flex items-center justify-center w-6 h-6 rounded-lg transition-all duration-300 bg-red-100 text-red-600"
				textClass = "text-xs font-semibold transition-all duration-300 text-red-600"
			} else {
				iconContainerClass = "flex items-center justify-center w-6 h-6 rounded-lg transition-all duration-300 group-hover:bg-red-50"
				textClass = "text-xs font-medium transition-all duration-300 text-red-600 group-hover:text-red-700"
			}
		} else if isActive {
			iconContainerClass = "flex items-center justify-center w-6 h-6 rounded-lg transition-all duration-300 bg-blue-100 text-blue-600"
			textClass = "text-xs font-semibold transition-all duration-300 text-blue-600"
		} else {
			iconContainerClass = "flex items-center justify-center w-6 h-6 rounded-lg transition-all duration-300 group-hover:bg-slate-100"
			textClass = "text-xs font-medium transition-all duration-300 text-slate-500 group-hover:text-slate-700"
		}

		return A(
			Href(href),
			Class("flex flex-col items-center justify-center p-3 transition-all duration-300 ease-out rounded-2xl group"),
			Iff(title == "Sign Out", func() Node {
				return Attr("hx-boost", "false") // Disable HTMX boost for reliable logout
			}),
			Div(
				Class("relative mb-1"),
				Div(
					Class(iconContainerClass),
					Div(
						Class("w-4 h-4"),
						icon,
					),
				),
			),
			Span(
				Class(textClass),
				Text(title),
			),
		)
	}

	return Nav(
		Class("fixed bottom-0 left-0 right-0 z-50 bg-white/95 backdrop-blur-lg border-t border-slate-200/50 px-4 py-2 shadow-2xl lg:hidden"),
		Div(
			Class("max-w-md mx-auto"),
			// Use different layout for unverified vs verified users
			If(r.AuthUser != nil && !r.AuthUser.Verified,
				// Closer spacing for unverified users with only 2 buttons
				Div(
					Class("flex items-center justify-center gap-16"),
					userNavItem(icons.Home(), "Home", "/home"),
					userNavItem(icons.Exit(), "Sign Out", r.Path(routenames.Logout)),
				),
			),
			If(r.AuthUser != nil && r.AuthUser.Verified,
				// Normal spacing for verified users with multiple buttons
				Div(
					Class("flex items-center justify-around"),
					// Show all nav items for verified users
					userNavItem(icons.Home(), "Home", "/home"),
					userNavItem(icons.CircleStack(), "Dashboard", "/dashboard"),
					userNavItem(icons.Star(), "Quizzes", "/quizzes"),
					userNavItem(icons.Archive(), "Notes", "/notes"),
					// More button for verified users only
					func() Node {
						// Check if current page is in the "More" section
						isMoreActive := r.CurrentPath == "/profile" || r.CurrentPath == "/files" ||
							strings.HasPrefix(r.CurrentPath, "/profile/") ||
							strings.HasPrefix(r.CurrentPath, "/admin") ||
							r.CurrentPath == "/cache"
						
						var iconContainerClass string
						var textClass string
						
						if isMoreActive {
							iconContainerClass = "flex items-center justify-center w-6 h-6 rounded-lg transition-all duration-300 bg-blue-100 text-blue-600"
							textClass = "text-xs font-semibold transition-all duration-300 text-blue-600"
						} else {
							iconContainerClass = "flex items-center justify-center w-6 h-6 rounded-lg transition-all duration-300 group-hover:bg-slate-100"
							textClass = "text-xs font-medium transition-all duration-300 text-slate-500 group-hover:text-slate-700"
						}
						
						return Button(
							ID("more_nav_button"),
							Class("flex flex-col items-center justify-center p-3 transition-all duration-300 ease-out rounded-2xl group"),
							Attr("onclick", "document.getElementById('mobile_nav_modal').classList.toggle('hidden'); this.classList.toggle('more-active');"),
							Div(
								Class("relative mb-1"),
								Div(
									Class(iconContainerClass),
									Div(
										Class("w-4 h-4"),
										Span(Class(func() string {
											if isMoreActive {
												return "text-blue-600 transition-colors duration-300"
											}
											return "text-slate-500 group-hover:text-slate-700 transition-colors duration-300"
										}()), Text("⋮⋮")),
									),
								),
							),
							Span(
								Class(textClass),
								Text("More"),
							),
						)
					}(),
				),
			),
		),
	)
}

func mobileNavModal(r *ui.Request) Node {
	return Div(
		ID("mobile_nav_modal"),
		Class("fixed bottom-20 right-4 bg-white rounded-lg shadow-lg border border-gray-200 z-50 hidden"),
		Style("min-width: 200px;"),
		Div(
			Class("p-4"),
			H3(
				Class("font-bold text-sm mb-3 text-gray-700"),
				Text("More Options"),
			),
			Div(
				Class("space-y-1"),
				A(
					Href("/profile"),
					Class("flex items-center gap-3 p-2 rounded-md hover:bg-slate-100 transition-colors text-sm"),
					Attr("onclick", "document.getElementById('mobile_nav_modal').classList.add('hidden'); document.getElementById('more_nav_button').classList.remove('more-active');"),
					icons.UserCircle(),
					Text("Profile"),
				),
				A(
					Href("/files"),
					Class("flex items-center gap-3 p-2 rounded-md hover:bg-slate-100 transition-colors text-sm"),
					Attr("onclick", "document.getElementById('mobile_nav_modal').classList.add('hidden'); document.getElementById('more_nav_button').classList.remove('more-active');"),
					icons.Document(),
					Text("Files"),
				),
				If(r.IsAdmin,
					Group{
						Hr(Class("my-2")),
						H4(Class("text-xs font-semibold text-slate-600 px-2 py-1"), Text("Admin Tools")),
						A(
							Href("/admin"),
							Class("flex items-center gap-3 p-2 rounded-md hover:bg-slate-100 transition-colors text-sm"),
							Attr("onclick", "document.getElementById('mobile_nav_modal').classList.add('hidden'); document.getElementById('more_nav_button').classList.remove('more-active');"),
							icons.Archive(),
							Text("Admin Panel"),
						),
						A(
							Href("/cache"),
							Class("flex items-center gap-3 p-2 rounded-md hover:bg-slate-100 transition-colors text-sm"),
							Attr("onclick", "document.getElementById('mobile_nav_modal').classList.add('hidden'); document.getElementById('more_nav_button').classList.remove('more-active');"),
							icons.Archive(),
							Text("Cache Management"),
						),
						A(
							Href("/admin/tasks"),
							Class("flex items-center gap-3 p-2 rounded-md hover:bg-slate-100 transition-colors text-sm"),
							Attr("onclick", "document.getElementById('mobile_nav_modal').classList.add('hidden'); document.getElementById('more_nav_button').classList.remove('more-active');"),
							icons.CircleStack(),
							Text("Background Tasks"),
						),
						A(
							Href("/admin/entity/user"),
							Class("flex items-center gap-3 p-2 rounded-md hover:bg-slate-100 transition-colors text-sm"),
							Attr("onclick", "document.getElementById('mobile_nav_modal').classList.add('hidden'); document.getElementById('more_nav_button').classList.remove('more-active');"),
							icons.UserCircle(),
							Text("User Management"),
						),
					},
				),
				// Sign Out placed at the bottom
				Hr(Class("my-2")),
				A(
					Href(r.Path(routenames.Logout)),
					Class("flex items-center gap-3 p-2 rounded-md hover:bg-red-50 text-red-600 transition-colors text-sm"),
					Attr("onclick", "document.getElementById('mobile_nav_modal').classList.add('hidden'); document.getElementById('more_nav_button').classList.remove('more-active');"),
					Attr("hx-boost", "false"), // Disable HTMX boost for reliable logout
					icons.Exit(),
					Text("Sign Out"),
				),
			),
		),
	)
}

func search() Node {
	return cache.SetIfNotExists("layout.search", func() Node {
		return Div(
			Class("ml-2"),
			Attr("x-data", ""),
			Label(
				Class("input"),
				icons.MagnifyingGlass(),
				Input(
					Type("search"),
					Class("grow"),
					Placeholder("Search"),
					Attr("@click", "search_modal.showModal();"),
				),
			),
		)
	})
}

func searchModal(r *ui.Request) Node {
	return cache.SetIfNotExists("layout.searchModal", func() Node {
		return Dialog(
			ID("search_modal"),
			Class("modal"),
			Div(
				Class("modal-box"),
				Form(
					Method("dialog"),
					Button(
						Class("btn btn-sm btn-circle btn-ghost absolute right-2 top-2"),
						Text("✕"),
					),
				),
				H3(
					Class("text-lg font-bold mb-2"),
					Text("Search"),
				),
				Input(
					Attr("hx-get", r.Path(routenames.Search)),
					Attr("hx-trigger", "keyup changed delay:500ms"),
					Attr("hx-target", "#results"),
					Name("query"),
					Class("input w-full"),
					Type("search"),
					Placeholder("Search..."),
				),
				Ul(
					ID("results"),
					Class("list"),
				),
			),
			Form(
				Method("dialog"),
				Class("modal-backdrop"),
				Button(
					Text("close"),
				),
			),
		)
	})
}

// unifiedHeader shows the same header for all users with different navigation based on auth status
func unifiedHeader(r *ui.Request) Node {
	// Basic safety check
	if r == nil {
		return Header(
			Class("bg-white border-b border-gray-200 shadow-sm"),
			Div(
				Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between"),
				Span(Class("text-2xl font-bold text-indigo-600"), Text("Zero")),
			),
		)
	}

	return Header(
		Class("bg-white border-b border-gray-200 shadow-sm sticky top-0 z-50"),
		Div(
			Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
			Div(
				Class("flex justify-between items-center h-16"),
				// Logo/Brand
				A(
					Href("/"),
					Class("flex items-center space-x-2"),
					Div(
						Class("w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center"),
						Span(Class("text-white font-bold text-lg"), Text("Z")),
					),
					Span(
						Class("text-xl font-bold text-gray-900"),
						Text("Zero"),
					),
				),
				// Navigation
				If(!r.IsAuth,
					// Non-authenticated navigation
					Group{
						// Desktop navigation with Sign In and Get Started
						Nav(
							Class("hidden md:flex items-center space-x-8"),
							Div(
								Class("flex items-center space-x-4"),
								A(
									Href("/user/login"),
									Class("text-gray-600 hover:text-emerald-600 font-medium transition-colors duration-200"),
									Text("Sign In"),
								),
								A(
									Href("/user/register"),
									Class("bg-emerald-600 hover:bg-emerald-700 text-white px-4 py-2 rounded-lg font-medium transition-colors duration-200"),
									Text("Sign Up"),
								),
							),
						),
						// Mobile navigation with just Sign In and Sign Up buttons
						Div(
							Class("flex md:hidden items-center space-x-3"),
							A(
								Href("/user/login"),
								Class("text-gray-600 hover:text-emerald-600 font-medium transition-colors duration-200 text-sm"),
								Text("Sign In"),
							),
							A(
								Href("/user/register"),
								Class("bg-emerald-600 hover:bg-emerald-700 text-white px-3 py-2 rounded-lg font-medium transition-colors duration-200 text-sm"),
								Text("Sign Up"),
							),
						),
					},
				),
				If(r.IsAuth && r.AuthUser != nil,
					// Authenticated navigation - user avatar and simple welcome
					Nav(
						Class("flex items-center space-x-4"),
						// User Avatar
						Div(
							Class("flex items-center space-x-3"),
							A(
							Href(func() string {
								if r != nil {
									return r.Path(routenames.Profile)
								}
								return "/profile"
							}()),
							Class("block"),
								func() Node {
								if r.AuthUser != nil && r.AuthUser.ProfilePicture != nil && *r.AuthUser.ProfilePicture != "" {
									return Img(
										Src("/files/"+*r.AuthUser.ProfilePicture),
										Alt("Profile Picture"),
										Class("w-8 h-8 rounded-full object-cover border-2 border-gray-200 hover:border-blue-300 transition-colors cursor-pointer"),
									)
								} else {
									return Div(
										Class("w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold text-sm hover:bg-blue-600 transition-colors cursor-pointer"),
										Text(func() string {
											if r.AuthUser != nil && r.AuthUser.Name != "" {
												return string(r.AuthUser.Name[0])
											}
											return "U"
										}()),
									)
								}
							}(),
							),
							// Welcome text (hidden on mobile)
							Div(
								Class("hidden md:block text-sm text-gray-600"),
								func() Node {
								name := "User"
								if r.AuthUser != nil && r.AuthUser.Name != "" {
									name = r.AuthUser.Name
								}
								return Span(
									Class("font-medium"),
									Text("Welcome, "+name),
								)
							}(),
							),
						),
					),
				),
			),
		),
	)
}

// unifiedFooter shows a simple footer for all users
func unifiedFooter(r *ui.Request) Node {
	// Basic safety check
	if r == nil {
		return Footer(
			Class("bg-gray-50 border-t border-gray-200 mt-16"),
			Div(
				Class("max-w-7xl mx-auto px-4 py-6 text-center"),
				P(Class("text-sm text-gray-600"), Text("© 2025 Zero. Your Quiz & Document Management Platform.")),
			),
		)
	}

	return Footer(
		Class("bg-gray-50 border-t border-gray-200 mt-16"),
		Div(
			Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8"),
			Div(
				Class("flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0"),
				// Brand section
				Div(
					Class("flex items-center space-x-2"),
					Div(
						Class("w-6 h-6 bg-indigo-600 rounded-lg flex items-center justify-center"),
						Span(Class("text-white font-bold text-sm"), Text("Z")),
					),
					Span(
						Class("text-lg font-bold text-gray-900"),
						Text("Zero"),
					),
				),
				// Simple links
				Div(
					Class("flex space-x-6 text-sm text-gray-600"),
					A(
						Href("/about"),
						Class("hover:text-indigo-600 transition-colors"),
						Text("About"),
					),
					A(
						Href("/contact"),
						Class("hover:text-indigo-600 transition-colors"),
						Text("Contact"),
					),
					If(!r.IsAuth,
						Group{
							A(
								Href("/user/login"),
								Class("hover:text-indigo-600 transition-colors"),
								Text("Sign In"),
							),
							A(
								Href("/user/register"),
								Class("hover:text-indigo-600 transition-colors"),
								Text("Sign Up"),
							),
						},
					),
					If(r.IsAuth,
						Group{
							A(
								Href("/profile"),
								Class("hover:text-indigo-600 transition-colors"),
								Text("Profile"),
							),
							A(
						Href(func() string {
							if r != nil {
								return r.Path(routenames.Logout)
							}
							return "/logout"
						}()),
						Class("hover:text-red-600 transition-colors"),
						Text("Sign Out"),
					),
						},
					),
				),
				// Copyright
				P(Class("text-sm text-gray-500"), Text("© 2025 Zero. All rights reserved.")),
			),
		),
	)
}
