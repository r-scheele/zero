package layouts

import (
	"strings"
	
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	"github.com/r-scheele/zero/pkg/ui/icons"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Admin provides a clean admin-only layout with responsive navigation.
// It includes a fixed header, sidebar navigation for desktop, and bottom navigation for mobile.
func Admin(r *ui.Request, content Node) Node {
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
				JS(),
			),
			Body(
				Class("min-h-screen bg-gradient-to-br from-slate-50 to-gray-100"),
				// Use custom admin header without logout button
				adminHeaderSimple(r),
				Div(
					Class("flex pt-16 min-h-screen"), // Account for fixed header
					adminSidebar(r),
			Main(
				Class("flex-1 bg-slate-50 lg:ml-80"),
						Div(
							Class("px-4 sm:px-6 lg:px-8 xl:px-12 py-6 lg:py-8 pb-20 lg:pb-8"), // Extra bottom padding on mobile for nav
							Div(
								Class("max-w-7xl mx-auto"),
								ID("main-content"),
								If(len(r.Title) > 0, H1(
									Class("text-3xl font-bold mb-8 text-slate-900"),
									Text(r.Title),
								)),
								FlashMessages(r),
								Div(
									Class("space-y-6"),
									content,
								),
							),
						),
					),
				),
				// Footer (placeholder for future implementation)
				// unifiedFooter(r),
				adminBottomNavigation(r), // Mobile navigation
				HtmxListeners(r),
			),
		),
	)
}

// adminHeaderSimple provides a simplified header for the admin layout without logout functionality.
func adminHeaderSimple(r *ui.Request) Node {
	// Basic safety check
	if r == nil {
		return Header(
			Class("fixed top-0 left-0 right-0 z-50 glass-header"),
			Div(
				Class("w-full px-4 sm:px-6 lg:px-8 xl:px-12 py-1 lg:py-2"),
				Div(
					Class("flex items-center justify-between h-8 lg:h-10 w-full"),
					Span(Class("text-2xl font-bold text-indigo-600"), Text("Zero")),
				),
			),
		)
	}

	return Header(
		Class("fixed top-0 left-0 right-0 z-50 glass-header"),
		Div(
			Class("w-full px-4 sm:px-6 lg:px-8 xl:px-12 py-1 lg:py-2"),
			Div(
				Class("flex items-center justify-between h-8 lg:h-10 w-full"),
				// Logo section
				Div(
					Class("flex items-center gap-3 sm:gap-4 lg:gap-5 flex-shrink-0"),
					// Logo
					A(
						Href("/home"),
						Class("flex items-center gap-2 sm:gap-3 lg:gap-4 text-base sm:text-lg lg:text-xl font-bold text-slate-800 hover:text-indigo-600 transition-all duration-300 ease-out"),
						Div(
							Class("w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center"),
							Span(Class("text-white font-bold text-lg"), Text("Z")),
						),
						Span(
							Class("text-2xl font-bold text-indigo-600"),
							Text("Zero"),
						),
					),
				),
				// Right side - just user info, no logout button
				Div(
					Class("flex items-center gap-2 sm:gap-3 lg:gap-4 flex-shrink-0"),
					Div(
						Class("hidden sm:flex items-center gap-2 text-sm text-slate-600"),
						If(r.AuthUser != nil && r.AuthUser.Name != "", Span(Text(r.AuthUser.Name))),
						Span(Class("badge badge-sm bg-indigo-100 text-indigo-800 border-indigo-200"), Text("Admin")),
					),
				),
			),
		),
	)
}

// adminHeader provides a full-featured header for the admin layout with user info and logout.
func adminHeader(r *ui.Request) Node {
	return Header(
		Class("fixed top-0 left-0 right-0 z-50 bg-white/95 backdrop-blur-lg border-b border-slate-200/50 shadow-elegant-sm"),
		Div(
			Class("container mx-auto max-w-none lg:max-w-6xl xl:max-w-7xl px-4 sm:px-6 lg:px-8 xl:px-12 py-3 lg:py-4"),
			Div(
				Class("flex items-center justify-between"),
				// Left side - Logo and Admin title
				Div(
					Class("flex items-center gap-4"),
					A(
						Href("/admin"),
						Class("flex items-center gap-3 text-slate-900 hover:text-blue-600 transition-colors duration-300"),
						Div(
							Class("flex items-center justify-center w-10 h-10 bg-gradient-to-r from-blue-600 to-blue-700 rounded-xl shadow-lg"),
							Span(
								Class("text-white font-bold text-xl"),
								Text("A"),
							),
						),
						Div(
							H1(
								Class("text-xl lg:text-2xl font-bold"),
								Text("Admin Panel"),
							),
							P(
								Class("text-sm text-slate-600 hidden sm:block"),
								Text("System Administration"),
							),
						),
					),
				),
				// Right side - User info and logout
				Div(
					Class("flex items-center gap-4"),
					Div(
						Class("hidden sm:flex items-center gap-2 text-sm text-slate-600"),
						icons.UserCircle(),
						If(r.AuthUser != nil && r.AuthUser.Name != "", Span(Text(r.AuthUser.Name))),
						Span(Class("badge badge-sm bg-blue-100 text-blue-800 border-blue-200"), Text("Admin")),
					),
					A(
						Href(r.Path(routenames.Logout)),
						Class("btn btn-outline btn-sm gap-2"),
						icons.Exit(),
						Span(Class("hidden sm:inline"), Text("Logout")),
					),
				),
			),
		),
	)
}

// adminMenuItem creates a navigation menu item for the admin sidebar
func adminMenuItem(r *ui.Request, icon Node, title, href string) Node {
	// Special handling for Users navigation - highlight when on any user page
	isActive := r.CurrentPath == href
	if href == "/admin/entity/user" && (r.CurrentPath == href || strings.HasPrefix(r.CurrentPath, "/admin/entity/user/")) {
		isActive = true
	}

	var linkClasses, iconClasses, textClasses string
	
	// Special styling for Sign Out button
	if title == "Sign Out" {
		if isActive {
			linkClasses = "admin-nav-item admin-nav-active flex items-center gap-3 px-4 py-3 rounded-xl bg-red-50 text-red-600 border border-red-200 transition-all duration-200"
			iconClasses = "w-5 h-5 text-red-600"
			textClasses = "font-semibold text-sm text-red-600"
		} else {
			linkClasses = "admin-nav-item flex items-center gap-3 px-4 py-3 rounded-xl text-red-600 hover:bg-red-50 hover:text-red-700 transition-all duration-200"
			iconClasses = "w-5 h-5 text-red-500 group-hover:text-red-600"
			textClasses = "font-medium text-sm text-red-600 group-hover:text-red-700"
		}
	} else if isActive {
		linkClasses = "admin-nav-item admin-nav-active flex items-center gap-3 px-4 py-3 rounded-xl bg-blue-100 text-blue-700 border border-blue-200 transition-all duration-200"
		iconClasses = "w-5 h-5 text-blue-600"
		textClasses = "font-semibold text-sm text-blue-700"
	} else {
		linkClasses = "admin-nav-item flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-slate-900 transition-all duration-200"
		iconClasses = "w-5 h-5 text-slate-500 group-hover:text-blue-600"
		textClasses = "font-medium text-sm group-hover:text-slate-900"
	}

	return Li(
		A(
			Href(href),
			Class(linkClasses+" group"),
			Div(
				Class(iconClasses),
				icon,
			),
			Span(
				Class(textClasses),
				Text(title),
			),
		),
	)
}

// adminSidebar creates the desktop sidebar navigation for the admin layout.
func adminSidebar(r *ui.Request) Node {
	return Aside(
		Class("hidden lg:block fixed left-0 top-16 bottom-0 w-80 bg-white border-r border-slate-200 shadow-sm flex-shrink-0 z-40"),
		Div(
			Class("p-6 h-full flex flex-col overflow-y-auto"),
			Nav(
				Class("flex-1 flex flex-col"),
				// Main navigation items
				Ul(
					Class("space-y-2"),
					HxBoost(),
					adminMenuItem(r, icons.Home(), "Overview", "/admin"),
					adminMenuItem(r, icons.UserCircle(), "Users", "/admin/entity/user"),
					adminMenuItem(r, icons.Document(), "Notes", "/admin/entity/note"),
					adminMenuItem(r, icons.Archive(), "Cache Management", "/admin/cache"),
					adminMenuItem(r, icons.CircleStack(), "Background Tasks", r.Path(routenames.AdminTasks)),
				),
				// Spacer to push Sign Out to bottom
				Div(
					Class("flex-grow"),
				),
				// Sign Out at the very bottom with visual separation
				Div(
					Class("mt-auto pt-4 border-t border-slate-200"),
					Ul(
						Class("space-y-2"),
						adminMenuItem(r, icons.Exit(), "Sign Out", r.Path(routenames.Logout)),
					),
				),
			),
		),
	)
}

// adminBottomNavItem creates a navigation item for the mobile bottom navigation
func adminBottomNavItem(r *ui.Request, icon Node, title, href string) Node {
	// Special handling for Users navigation - highlight when on any user page
	isActive := r.CurrentPath == href
	if href == "/admin/entity/user" && (r.CurrentPath == href || strings.HasPrefix(r.CurrentPath, "/admin/entity/user/")) {
		isActive = true
	}

	var iconContainerClass, textClass string
	
	// Special styling for Logout button
	if title == "Logout" {
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
		Class("flex flex-col items-center justify-center p-2 transition-all duration-300 ease-out rounded-xl group"),
		Div(
			Class("relative mb-1"),
			Div(
				Class(iconContainerClass),
				Div(
				Class("w-4 h-4"),
				If(title == "Logout" && isActive,
					Span(Class("text-red-600"), icon),
				),
				If(title == "Logout" && !isActive,
					Span(Class("text-red-600 group-hover:text-red-700 transition-colors duration-300"), icon),
				),
				If(title != "Logout" && isActive,
					Span(Class("text-blue-600"), icon),
				),
				If(title != "Logout" && !isActive,
					Span(Class("text-slate-500 group-hover:text-blue-600 transition-colors duration-300"), icon),
				),
			),
			),
			If(isActive,
				Div(
					Class("absolute -bottom-1 left-1/2 transform -translate-x-1/2 w-1 h-1 bg-white rounded-full"),
				),
			),
		),
		Span(
			Class(textClass),
			Text(title),
		),
	)
}

// adminBottomNavigation creates the mobile bottom navigation for the admin layout.
func adminBottomNavigation(r *ui.Request) Node {

	return Nav(
		Class("fixed bottom-0 left-0 right-0 z-50 bg-white/95 backdrop-blur-lg border-t border-slate-200/50 px-3 py-3 shadow-2xl lg:hidden"),
		Div(
			Class("max-w-md mx-auto"),
			Div(
				Class("flex items-center justify-around"),
				adminBottomNavItem(r, icons.Home(), "Overview", "/admin"),
			adminBottomNavItem(r, icons.UserCircle(), "Users", "/admin/entity/user"),
			adminBottomNavItem(r, icons.Archive(), "Cache", r.Path(routenames.Cache)),
			adminBottomNavItem(r, icons.CircleStack(), "Tasks", r.Path(routenames.AdminTasks)),
			adminBottomNavItem(r, icons.Exit(), "Logout", r.Path(routenames.Logout)),
			),
		),
	)
}
