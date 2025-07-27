package pages

import (
	"fmt"
	"net/url"

	"entgo.io/ent/entc/load"
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/ent/admin"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AdminEntityDelete(ctx echo.Context, entityTypeName string) error {
	r := ui.NewRequest(ctx)
	r.Title = ""

	return r.Render(
		layouts.Admin,
		forms.AdminEntityDelete(r, entityTypeName),
	)
}

func AdminEntityView(ctx echo.Context, entityTypeName string, entity map[string][]string, id int) error {
	r := ui.NewRequest(ctx)
	r.Title = ""

	// Special handling for User entity
	if entityTypeName == "User" {
		return AdminUserView(ctx, entity, id)
	}

	// Generic entity view for other entities
	return r.Render(
		layouts.Admin,
		Div(
			Class("space-y-6"),
			Div(
				Class("bg-white rounded-lg shadow-sm border border-gray-200 p-6"),
				Div(
					Class("grid grid-cols-1 md:grid-cols-2 gap-4"),
					Group(func() []Node {
						var fields []Node
						fields = append(fields,
							Div(
								Class("border-b pb-2 mb-2"),
								Dt(Class("text-sm font-medium text-gray-500"), Text("ID")),
								Dd(Class("text-sm text-gray-900"), Text(fmt.Sprint(id))),
							),
						)
						for key, values := range entity {
							if len(values) > 0 {
								fields = append(fields,
									Div(
										Class("border-b pb-2 mb-2"),
										Dt(Class("text-sm font-medium text-gray-500"), Text(key)),
										Dd(Class("text-sm text-gray-900"), Text(values[0])),
									),
								)
							}
						}
						return fields
					}()),
				),
				Div(
					Class("flex gap-4 mt-6"),
					ButtonLink(
						ColorInfo,
						r.Path(routenames.AdminEntityEdit(entityTypeName), id),
						"Edit",
					),
					ButtonLink(
						ColorError,
						r.Path(routenames.AdminEntityDelete(entityTypeName), id),
						"Delete",
					),
				),
			),
		),
	)
}

func AdminEntityInput(ctx echo.Context, schema *load.Schema, values url.Values) error {
	r := ui.NewRequest(ctx)
	r.Title = ""

	return r.Render(
		layouts.Admin,
		forms.AdminEntity(r, schema, values),
	)
}

func AdminEntityList(
	ctx echo.Context,
	entityTypeName string,
	entityList *admin.EntityList,
) error {
	r := ui.NewRequest(ctx)
	r.Title = ""

	genHeader := func() Node {
		g := make(Group, 0, len(entityList.Columns)+1)
		g = append(g, Th(Class("sticky left-0 bg-white z-10 w-16"), Text("ID")))

		if entityTypeName == "User" {
			// For User entity, implement responsive column visibility with fixed widths
			responsiveColumns := []struct {
				name    string
				classes string
			}{
				{"Name", "w-32"}, // Always visible
				{"Phone number", "hidden sm:table-cell w-32"},      // Hidden on mobile
				{"Email", "hidden md:table-cell w-48"},             // Hidden on mobile and tablet
				{"Verified", "hidden lg:table-cell w-20"},          // Hidden except on large screens
				{"Verification code", "hidden xl:table-cell w-24"}, // Hidden except on extra large screens
				{"Admin", "w-20"}, // Always visible
				{"Registration method", "hidden xl:table-cell w-24"}, // Hidden except on extra large screens
				{"Profile picture", "hidden w-0"},                    // Always hidden
				{"Dark mode", "hidden w-0"},                          // Always hidden
				{"Bio", "hidden w-0"},                                // Always hidden
				{"Email notifications", "hidden w-0"},                // Always hidden
				{"Sms notifications", "hidden w-0"},                  // Always hidden
				{"Is active", "hidden lg:table-cell w-20"},           // Hidden except on large screens
				{"Last login", "hidden xl:table-cell w-32"},          // Hidden except on extra large screens
				{"Created at", "hidden lg:table-cell w-32"},          // Hidden except on large screens
				{"Updated at", "hidden w-0"},                         // Always hidden
			}

			for _, col := range responsiveColumns {
				g = append(g, Th(Class(col.classes), Text(col.name)))
			}
		} else {
			// For other entities, keep original behavior
			for _, h := range entityList.Columns {
				g = append(g, Th(Text(h)))
			}
			g = append(g, Th())
		}
		return g
	}

	genRow := func(row admin.EntityValues) Node {
		if entityTypeName == "User" {
			// For User entity, make entire row clickable with responsive column visibility
			g := make(Group, 0, len(row.Values)+1)
			g = append(g, Th(Class("sticky left-0 bg-white z-10 w-16"), Text(fmt.Sprint(row.ID))))

			// Apply responsive classes to match header with fixed widths
			responsiveClasses := []string{
				"w-32",                      // Name - Always visible
				"hidden sm:table-cell w-32", // Phone number - Hidden on mobile
				"hidden md:table-cell w-48", // Email - Hidden on mobile and tablet
				"hidden lg:table-cell w-20", // Verified - Hidden except on large screens
				"hidden xl:table-cell w-24", // Verification code - Hidden except on extra large screens
				"w-20",                      // Admin - Always visible
				"hidden xl:table-cell w-24", // Registration method - Hidden except on extra large screens
				"hidden w-0",                // Profile picture - Always hidden
				"hidden w-0",                // Dark mode - Always hidden
				"hidden w-0",                // Bio - Always hidden
				"hidden w-0",                // Email notifications - Always hidden
				"hidden w-0",                // Sms notifications - Always hidden
				"hidden lg:table-cell w-20", // Is active - Hidden except on large screens
				"hidden xl:table-cell w-32", // Last login - Hidden except on extra large screens
				"hidden lg:table-cell w-32", // Created at - Hidden except on large screens
				"hidden w-0",                // Updated at - Always hidden
			}

			for i, h := range row.Values {
				if i < len(responsiveClasses) {
					g = append(g, Td(Class(responsiveClasses[i]), Text(h)))
				} else {
					g = append(g, Td(Text(h)))
				}
			}

			return Tr(
				Class("cursor-pointer hover:bg-blue-50 transition-colors"),
				Attr("hx-get", r.Path(routenames.AdminEntityView(entityTypeName), row.ID)),
				Attr("hx-push-url", "true"),
				Attr("hx-target", "#main-content"),
				g,
			)
		} else {
			// For other entities, keep the old Edit/Delete buttons
			g := make(Group, 0, len(row.Values)+3)
			g = append(g, Th(Text(fmt.Sprint(row.ID))))
			for _, h := range row.Values {
				g = append(g, Td(Text(h)))
			}
			g = append(g,
				Td(
					ButtonLink(
						ColorInfo,
						r.Path(routenames.AdminEntityEdit(entityTypeName), row.ID),
						"Edit",
					),
					Span(Class("mr-2")),
					ButtonLink(
						ColorError,
						r.Path(routenames.AdminEntityDelete(entityTypeName), row.ID),
						"Delete",
					),
				),
			)
			return Tr(g)
		}
	}

	genRows := func() Node {
		g := make(Group, 0, len(entityList.Entities))
		for _, row := range entityList.Entities {
			g = append(g, genRow(row))
		}
		return g
	}

	// Search form for User entity
	searchForm := func() Node {
		if entityTypeName != "User" {
			return Group{}
		}

		searchValue := ctx.QueryParam("search")

		return Div(
			Class("mb-6"),
			Div(
				Class("flex gap-4 items-center"),
				Div(
					Class("flex-1"),
					Input(
						Type("text"),
						Name("search"),
						Placeholder("Search users by name, phone number, or email..."),
						Value(searchValue),
						Class("input input-bordered w-full"),
						Attr("hx-get", r.Path(routenames.AdminEntityList(entityTypeName))),
						Attr("hx-trigger", "input changed delay:300ms, search"),
						Attr("hx-target", "#user-table-container"),
						Attr("hx-include", "this"),
						Attr("hx-indicator", "#search-indicator"),
					),
				),
				Div(
					ID("search-indicator"),
					Class("htmx-indicator"),
					Div(Class("loading loading-spinner loading-sm")),
				),
				If(len(searchValue) > 0,
					A(
						Href(r.Path(routenames.AdminEntityList(entityTypeName))),
						Class("btn btn-ghost"),
						Text("Clear"),
					),
				),
			),
		)
	}

	content := Group{
		searchForm(),
		// Only show Add button for non-User entities
		If(entityTypeName != "User",
			Div(
				Class("form-control mb-2"),
				ButtonLink(
					ColorAccent,
					r.Path(routenames.AdminEntityAdd(entityTypeName)),
					fmt.Sprintf("Add %s", entityTypeName),
				),
			),
		),
		Div(
			ID("user-table-container"),
			Class("w-full"),
			Table(
				Class("table table-zebra mb-2 w-full"),
				THead(
					Tr(genHeader()),
				),
				TBody(genRows()),
			),
			Pager(
				entityList.Page,
				r.Path(routenames.AdminEntityList(entityTypeName)),
				entityList.HasNextPage,
				"",
			),
		),
	}

	return r.Render(layouts.Admin, content)
}

func AdminUserView(ctx echo.Context, entity map[string][]string, id int) error {
	r := ui.NewRequest(ctx)
	r.Title = ""

	getValue := func(key string) string {
		if values, ok := entity[key]; ok && len(values) > 0 {
			return values[0]
		}
		return "-"
	}

	isAdmin := getValue("admin") == "true"
	isVerified := getValue("verified") == "true"
	isActive := getValue("is_active") == "true"

	return r.Render(
		layouts.Admin,
		Div(
			Class("max-w-4xl mx-auto p-6 space-y-6"),

			// Header
			Div(
				Class("bg-white rounded-xl shadow-sm border border-gray-200 p-6"),
				Div(
					Class("flex items-start gap-6"),
					Div(
						Class("w-20 h-20 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white text-2xl font-bold"),
						Text(func() string {
							name := getValue("name")
							if len(name) > 0 {
								return string(name[0])
							}
							return "?"
						}()),
					),
					Div(
						Class("flex-1"),
						H1(Class("text-2xl font-bold text-gray-900 mb-2"), Text(getValue("name"))),
						Div(
							Class("flex flex-wrap gap-2 mb-3"),
							Span(
								Class("px-3 py-1 text-xs font-medium rounded-full bg-blue-100 text-blue-800"),
								Text(fmt.Sprintf("ID: %d", id)),
							),
							If(isAdmin,
								Span(Class("px-3 py-1 text-xs font-medium rounded-full bg-purple-100 text-purple-800"), Text("Admin")),
							),
							If(isVerified,
								Span(Class("px-3 py-1 text-xs font-medium rounded-full bg-green-100 text-green-800"), Text("Verified")),
							),
							If(!isVerified,
								Span(Class("px-3 py-1 text-xs font-medium rounded-full bg-red-100 text-red-800"), Text("Unverified")),
							),
							If(isActive,
								Span(Class("px-3 py-1 text-xs font-medium rounded-full bg-emerald-100 text-emerald-800"), Text("Active")),
							),
							If(!isActive,
								Span(Class("px-3 py-1 text-xs font-medium rounded-full bg-gray-100 text-gray-800"), Text("Inactive")),
							),
						),
						P(Class("text-gray-600"), Text(getValue("email"))),
					),
				),
			),

			// Personal Info
			Div(
				Class("bg-white rounded-xl shadow-sm border border-gray-200 p-6"),
				H2(Class("text-lg font-semibold text-gray-900 mb-4"), Text("Personal Information")),
				Div(
					Class("grid grid-cols-1 md:grid-cols-2 gap-6"),
					Div(
						Class("space-y-4"),
						Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Full Name")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("name")))),
						Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Email")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("email")))),
						Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Phone Number")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("phone_number")))),
					),
					Div(
						Class("space-y-4"),
						Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Registration Method")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("registration_method")))),
						Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Verification Code")), Dd(Class("mt-1 text-sm text-gray-900 font-mono"), Text(getValue("verification_code")))),
						Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Bio")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("bio")))),
					),
				),
			),

			// Account Info
			Div(
				Class("bg-white rounded-xl shadow-sm border border-gray-200 p-6"),
				H2(Class("text-lg font-semibold text-gray-900 mb-4"), Text("Account Information")),
				Div(
					Class("grid grid-cols-1 md:grid-cols-2 gap-6"),
					Div(
						Class("space-y-4"),
						Div(
							Dt(Class("text-sm font-medium text-gray-500"), Text("Account Status")),
							Dd(Class("mt-1"),
								If(isActive,
									Span(Class("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"), Text("Active")),
								),
								If(!isActive,
									Span(Class("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800"), Text("Inactive")),
								),
							),
						),
						Div(
							Dt(Class("text-sm font-medium text-gray-500"), Text("Email Verified")),
							Dd(Class("mt-1"),
								If(isVerified,
									Span(Class("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"), Text("Verified")),
								),
								If(!isVerified,
									Span(Class("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800"), Text("Pending")),
								),
							),
						),
						Div(
							Dt(Class("text-sm font-medium text-gray-500"), Text("Admin Status")),
							Dd(Class("mt-1"),
								If(isAdmin,
									Span(Class("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800"), Text("Administrator")),
								),
								If(!isAdmin,
									Span(Class("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800"), Text("Regular User")),
								),
							),
						),
					),
					Div(
						Class("space-y-4"),
						Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Dark Mode")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("dark_mode")))),
						Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Email Notifications")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("email_notifications")))),
						Div(Dt(Class("text-sm font-medium text-gray-500"), Text("SMS Notifications")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("sms_notifications")))),
					),
				),
			),

			// Activity Info
			Div(
				Class("bg-white rounded-xl shadow-sm border border-gray-200 p-6"),
				H2(Class("text-lg font-semibold text-gray-900 mb-4"), Text("Activity Information")),
				Div(
					Class("grid grid-cols-1 md:grid-cols-2 gap-6"),
					Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Last Login")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("last_login")))),
					Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Account Created")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("created_at")))),
					Div(Dt(Class("text-sm font-medium text-gray-500"), Text("Last Updated")), Dd(Class("mt-1 text-sm text-gray-900"), Text(getValue("updated_at")))),
				),
			),

			// Actions
			Div(
				Class("bg-white rounded-xl shadow-sm border border-gray-200 p-6"),
				H2(Class("text-lg font-semibold text-gray-900 mb-4"), Text("Account Management")),
				Div(
					Class("flex flex-wrap gap-3"),

					// Profile Actions
					Div(
						Class("flex gap-2"),
						H3(Class("text-sm font-medium text-gray-700 mb-2 w-full"), Text("Profile Actions")),
						A(
							Class("inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"),
							Href(r.Path(routenames.AdminEntityEdit("User"), id)),
							Text("✏️ Edit Profile"),
						),
					),

					// Communication
					Div(
						Class("flex gap-2"),
						H3(Class("text-sm font-medium text-gray-700 mb-2 w-full"), Text("Communication")),
						Button(
							Class("inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"),
							Attr("onclick", fmt.Sprintf("window.open('mailto:%s', '_blank')", getValue("Email"))),
							Text("📧 Send Email"),
						),
						If(getValue("Phone") != "-",
							Button(
								Class("inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"),
								Attr("onclick", fmt.Sprintf("window.open('sms:%s', '_blank')", getValue("Phone"))),
								Text("📱 Send SMS"),
							),
						),
					),

					// Danger Zone
					Div(
						Class("w-full border-t pt-4 mt-4"),
						H3(Class("text-sm font-medium text-red-700 mb-2"), Text("⚠️ Danger Zone")),
						A(
							Class("inline-flex items-center px-4 py-2 border border-red-300 rounded-md shadow-sm text-sm font-medium text-red-700 bg-white hover:bg-red-50"),
							Href(r.Path(routenames.AdminEntityDelete("User"), id)),
							Text("🗑️ Delete Account"),
						),
					),
				),
			),
		),
	)
}

func AdminEntityListTable(
	ctx echo.Context,
	entityTypeName string,
	entityList *admin.EntityList,
) error {
	r := ui.NewRequest(ctx)

	genHeader := func() Node {
		g := make(Group, 0, len(entityList.Columns)+1)
		g = append(g, Th(Class("sticky left-0 bg-white z-10"), Text("ID")))

		if entityTypeName == "User" {
			// For User entity, show minimal columns to prevent horizontal scrolling
			responsiveColumns := []struct {
				name    string
				classes string
			}{
				{"Name", ""},                            // Always visible
				{"Phone number", "hidden"},              // Always hidden
				{"Email", "hidden xl:table-cell"},       // Hidden except on extra large screens
				{"Verified", "hidden 2xl:table-cell"},   // Hidden except on 2xl screens
				{"Verification code", "hidden"},         // Always hidden
				{"Admin", ""},                           // Always visible
				{"Registration method", "hidden"},       // Always hidden
				{"Profile picture", "hidden"},           // Always hidden
				{"Dark mode", "hidden"},                 // Always hidden
				{"Bio", "hidden"},                       // Always hidden
				{"Email notifications", "hidden"},       // Always hidden
				{"Sms notifications", "hidden"},         // Always hidden
				{"Is active", "hidden 2xl:table-cell"},  // Hidden except on 2xl screens
				{"Last login", "hidden"},                // Always hidden
				{"Created at", "hidden 2xl:table-cell"}, // Hidden except on 2xl screens
				{"Updated at", "hidden"},                // Always hidden
			}

			for _, col := range responsiveColumns {
				g = append(g, Th(Class(col.classes), Text(col.name)))
			}
		} else {
			// For other entities, keep original behavior
			for _, h := range entityList.Columns {
				g = append(g, Th(Text(h)))
			}
			g = append(g, Th())
		}
		return g
	}

	genRow := func(row admin.EntityValues) Node {
		if entityTypeName == "User" {
			// For User entity, make entire row clickable with responsive column visibility
			g := make(Group, 0, len(row.Values)+1)
			g = append(g, Th(Class("sticky left-0 bg-white z-10"), Text(fmt.Sprint(row.ID))))

			// Apply responsive classes to match header - minimal columns only
			responsiveClasses := []string{
				"",                      // Name - Always visible
				"hidden",                // Phone number - Always hidden
				"hidden xl:table-cell",  // Email - Hidden except on extra large screens
				"hidden 2xl:table-cell", // Verified - Hidden except on 2xl screens
				"hidden",                // Verification code - Always hidden
				"",                      // Admin - Always visible
				"hidden",                // Registration method - Always hidden
				"hidden",                // Profile picture - Always hidden
				"hidden",                // Dark mode - Always hidden
				"hidden",                // Bio - Always hidden
				"hidden",                // Email notifications - Always hidden
				"hidden",                // Sms notifications - Always hidden
				"hidden 2xl:table-cell", // Is active - Hidden except on 2xl screens
				"hidden",                // Last login - Always hidden
				"hidden 2xl:table-cell", // Created at - Hidden except on 2xl screens
				"hidden",                // Updated at - Always hidden
			}

			for i, h := range row.Values {
				if i < len(responsiveClasses) {
					g = append(g, Td(Class(responsiveClasses[i]), Text(h)))
				} else {
					g = append(g, Td(Text(h)))
				}
			}

			return Tr(
				Class("cursor-pointer hover:bg-blue-50 transition-colors"),
				Attr("hx-get", r.Path(routenames.AdminEntityView(entityTypeName), row.ID)),
				Attr("hx-push-url", "true"),
				Attr("hx-target", "#main-content"),
				g,
			)
		} else {
			// For other entities, keep the old Edit/Delete buttons
			g := make(Group, 0, len(row.Values)+3)
			g = append(g, Th(Text(fmt.Sprint(row.ID))))
			for _, h := range row.Values {
				g = append(g, Td(Text(h)))
			}
			g = append(g,
				Td(
					ButtonLink(
						ColorInfo,
						r.Path(routenames.AdminEntityEdit(entityTypeName), row.ID),
						"Edit",
					),
					Span(Class("mr-2")),
					ButtonLink(
						ColorError,
						r.Path(routenames.AdminEntityDelete(entityTypeName), row.ID),
						"Delete",
					),
				),
			)
			return Tr(g)
		}
	}

	genRows := func() Node {
		g := make(Group, 0, len(entityList.Entities))
		for _, row := range entityList.Entities {
			g = append(g, genRow(row))
		}
		return g
	}

	// Return only the table content for HTMX requests
	content := Div(
		Table(
			Class("table table-zebra mb-2 w-full"),
			THead(
				Tr(genHeader()),
			),
			TBody(genRows()),
		),
		Pager(
			entityList.Page,
			r.Path(routenames.AdminEntityList(entityTypeName)),
			entityList.HasNextPage,
			"",
		),
	)

	return content.Render(r.Context.Response().Writer)
}
