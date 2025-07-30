package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/backlite/ui"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/ent/note"
	"github.com/r-scheele/zero/ent/user"
	"github.com/r-scheele/zero/pkg/context"
	"github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Admin struct {
	orm      *ent.Client
	graph    *gen.Graph
	backlite *ui.Handler
	auth     *services.AuthClient
}

func init() {
	Register(new(Admin))
}

func (h *Admin) Init(c *services.Container) error {
	var err error
	h.graph = c.Graph
	h.orm = c.ORM
	h.auth = c.Auth
	h.backlite, err = ui.NewHandler(ui.Config{
		DB:           c.Database,
		BasePath:     "/admin/tasks",
		ItemsPerPage: 25,
		ReleaseAfter: c.Config.Tasks.ReleaseAfter,
	})
	return err
}

func (h *Admin) Routes(g *echo.Group) {
	ag := g.Group("/admin", middleware.RequireAdmin)

	// Admin overview/dashboard
	ag.GET("", h.Overview()).Name = "admin:overview"

	entities := ag.Group("/entity")
	for _, n := range h.graph.Nodes {
		// Skip PasswordToken entity for security reasons
		if n.Name == "PasswordToken" {
			continue
		}

		ng := entities.Group(fmt.Sprintf("/%s", strings.ToLower(n.Name)))
		ng.GET("", h.EntityList(n)).
			Name = routenames.AdminEntityList(n.Name)
		ng.GET("/add", h.EntityAdd(n)).
			Name = routenames.AdminEntityAdd(n.Name)
		ng.POST("/add", h.EntityAddSubmit(n)).
			Name = routenames.AdminEntityAddSubmit(n.Name)
		ng.GET("/:id", h.EntityView(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityView(n.Name)
		ng.GET("/:id/edit", h.EntityEdit(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityEdit(n.Name)
		ng.POST("/:id/edit", h.EntityEditSubmit(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityEditSubmit(n.Name)
		ng.GET("/:id/delete", h.EntityDelete(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityDelete(n.Name)
		ng.POST("/:id/delete", h.EntityDeleteSubmit(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityDeleteSubmit(n.Name)
	}

	// User-specific admin actions
	userGroup := ag.Group("/user")
	userGroup.POST("/:id/verify", h.VerifyUser)

	tasks := ag.Group("/tasks")
	tasks.GET("", h.AdminTasks).Name = routenames.AdminTasks
	tasks.GET("/succeeded", h.Backlite(h.backlite.Succeeded))
	tasks.GET("/failed", h.Backlite(h.backlite.Failed))
	tasks.GET("/upcoming", h.Backlite(h.backlite.Upcoming))
	tasks.GET("/task/:id", h.Backlite(h.backlite.Task))
	tasks.GET("/completed/:id", h.Backlite(h.backlite.TaskCompleted))
}

// Overview displays the admin dashboard
func (h *Admin) Overview() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return pages.AdminOverview(ctx, h.orm)
	}
}

// AdminTasks displays the admin tasks page with proper layout
func (h *Admin) AdminTasks(ctx echo.Context) error {
	return pages.AdminTasks(ctx)
}

// middlewareEntityLoad is middleware to extract the entity ID and attempt to load the given entity.
func (h *Admin) middlewareEntityLoad(n *gen.Type) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			id, err := strconv.Atoi(ctx.Param("id"))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid entity ID")
			}

			// Store the entity ID in context
			ctx.Set(context.AdminEntityIDKey, id)

			// Load entity data based on type
			switch n.Name {
			case "User":
				user, err := h.orm.User.Get(ctx.Request().Context(), id)
				if err != nil {
					return echo.NewHTTPError(http.StatusNotFound, "user not found")
				}
				
				// Handle nullable fields
				email := ""
				if user.Email != nil {
					email = *user.Email
				}
				
				verificationCode := ""
				if user.VerificationCode != nil {
					verificationCode = *user.VerificationCode
				}
				
				profilePicture := ""
				if user.ProfilePicture != nil {
					profilePicture = *user.ProfilePicture
				}
				
				bio := ""
				if user.Bio != nil {
					bio = *user.Bio
				}
				
				lastLogin := ""
				if user.LastLogin != nil {
					lastLogin = user.LastLogin.Format("2006-01-02 15:04")
				}
				
				updatedAt := ""
				if user.UpdatedAt != nil {
					updatedAt = user.UpdatedAt.Format("2006-01-02 15:04")
				}
				
				// Populate all user fields for the detail view
				entityData := map[string][]string{
					"id":                   {strconv.Itoa(user.ID)},
					"name":                 {user.Name},
					"phone_number":         {user.PhoneNumber},
					"email":                {email},
					"verified":             {strconv.FormatBool(user.Verified)},
					"verification_code":    {verificationCode},
					"admin":                {strconv.FormatBool(user.Admin)},
					"registration_method":  {string(user.RegistrationMethod)},
					"profile_picture":      {profilePicture},
					"dark_mode":            {strconv.FormatBool(user.DarkMode)},
					"bio":                  {bio},
					"email_notifications":  {strconv.FormatBool(user.EmailNotifications)},
					"sms_notifications":    {strconv.FormatBool(user.SmsNotifications)},
					"is_active":            {strconv.FormatBool(user.IsActive)},
					"last_login":           {lastLogin},
					"created_at":           {user.CreatedAt.Format("2006-01-02 15:04")},
					"updated_at":           {updatedAt},
				}
				ctx.Set(context.AdminEntityKey, entityData)
			case "Note":
				note, err := h.orm.Note.Query().Where(note.ID(id)).WithOwner().Only(ctx.Request().Context())
				if err != nil {
					return echo.NewHTTPError(http.StatusNotFound, "note not found")
				}
				
				ownerName := "Unknown"
				if note.Edges.Owner != nil {
					ownerName = note.Edges.Owner.Name
				}
				
				description := note.Description
				content := note.Content
				aiCurriculum := note.AiCurriculum
				shareToken := note.ShareToken
				
				entityData := map[string][]string{
					"id":                {strconv.Itoa(note.ID)},
					"title":             {note.Title},
					"description":       {description},
					"content":           {content},
					"visibility":        {string(note.Visibility)},
					"permission_level":  {string(note.PermissionLevel)},
					"share_token":       {shareToken},
					"ai_processing":     {strconv.FormatBool(note.AiProcessing)},
					"ai_curriculum":     {aiCurriculum},
					"owner":             {ownerName},
					"created_at":        {note.CreatedAt.Format("2006-01-02 15:04")},
					"updated_at":        {note.UpdatedAt.Format("2006-01-02 15:04")},
				}
				ctx.Set(context.AdminEntityKey, entityData)
			default:
				return echo.NewHTTPError(http.StatusNotImplemented, "entity type not supported")
			}

			return next(ctx)
		}
	}
}

func (h *Admin) EntityList(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		switch n.Name {
		case "User":
			// Get search query parameter
			searchQuery := ctx.QueryParam("search")
			
			// Build user query with optional search filtering
			userQuery := h.orm.User.Query()
			if searchQuery != "" {
				// Search in name, phone number, and email fields
				userQuery = userQuery.Where(
					user.Or(
						user.NameContains(searchQuery),
						user.PhoneNumberContains(searchQuery),
						user.EmailContains(searchQuery),
					),
				)
			}
			
			users, err := userQuery.All(ctx.Request().Context())
			if err != nil {
				fmt.Printf("Error fetching users: %v\n", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch users: "+err.Error())
			}
			
			fmt.Printf("Found %d users in database\n", len(users))
			
			entityList := &pages.EntityList{
				Columns: []string{
					"Name", "Phone number", "Email", "Verified", "Verification code",
					"Admin", "Registration method", "Profile picture", "Dark mode", "Bio",
					"Email notifications", "Sms notifications", "Is active", "Last login",
					"Created at", "Updated at",
				},
				Entities: make([]pages.EntityValues, len(users)),
			}
			
			for i, user := range users {
				email := ""
				if user.Email != nil {
					email = *user.Email
				}
				
				verificationCode := ""
				if user.VerificationCode != nil {
					verificationCode = *user.VerificationCode
				}
				
				profilePicture := ""
				if user.ProfilePicture != nil {
					profilePicture = *user.ProfilePicture
				}
				
				bio := ""
				if user.Bio != nil {
					bio = *user.Bio
				}
				
				lastLogin := ""
				if user.LastLogin != nil {
					lastLogin = user.LastLogin.Format("2006-01-02 15:04")
				}
				
				updatedAt := ""
				if user.UpdatedAt != nil {
					updatedAt = user.UpdatedAt.Format("2006-01-02 15:04")
				}
				
				fmt.Printf("User %d: ID=%d, Name=%s, Email=%s\n", i, user.ID, user.Name, email)
				
				// Provide all 16 values to match the UI columns with actual data
				entityList.Entities[i] = pages.EntityValues{
					ID: user.ID,
					Values: []string{
						user.Name,                                    // Name
						user.PhoneNumber,                            // Phone number
						email,                                       // Email
						strconv.FormatBool(user.Verified),          // Verified
						verificationCode,                           // Verification code
						strconv.FormatBool(user.Admin),             // Admin
						string(user.RegistrationMethod),            // Registration method
						profilePicture,                             // Profile picture
						strconv.FormatBool(user.DarkMode),          // Dark mode
						bio,                                        // Bio
						strconv.FormatBool(user.EmailNotifications), // Email notifications
						strconv.FormatBool(user.SmsNotifications),  // SMS notifications
						strconv.FormatBool(user.IsActive),          // Is active
						lastLogin,                                  // Last login
						user.CreatedAt.Format("2006-01-02 15:04"), // Created at
						updatedAt,                                  // Updated at
					},
				}
			}
			
			// Check if this is an HTMX request specifically for search functionality
			// Only return table content if it's an HTMX request with a search parameter or targeting the table container
			hxTarget := ctx.Request().Header.Get("HX-Target")
			if ctx.Request().Header.Get("HX-Request") == "true" && 
				(ctx.QueryParam("search") != "" || hxTarget == "#note-table-container" || hxTarget == "#user-table-container") {
				// Return only the table content for HTMX search requests
				return pages.AdminEntityListTable(ctx, n.Name, entityList)
			}
			
			return pages.AdminEntityList(ctx, n.Name, entityList)
		case "Note":
			// Get search query parameter
			searchQuery := ctx.QueryParam("search")
			
			// Build note query with user relationship and optional search filtering
			noteQuery := h.orm.Note.Query().WithOwner()
			if searchQuery != "" {
				// Search in title, content, and owner name
				noteQuery = noteQuery.Where(
					note.Or(
						note.TitleContains(searchQuery),
						note.ContentContains(searchQuery),
						note.HasOwnerWith(user.NameContains(searchQuery)),
					),
				)
			}
			
			notes, err := noteQuery.All(ctx.Request().Context())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch notes")
			}
			
			entityList := &pages.EntityList{
				Columns:  []string{"Title", "Content", "Created By", "Visibility", "Created At"},
				Entities: make([]pages.EntityValues, len(notes)),
			}
			
			for i, note := range notes {
				content := note.Content
				if len(content) > 100 {
					content = content[:100] + "..."
				}
				
				ownerName := "Unknown"
				if note.Edges.Owner != nil {
					ownerName = note.Edges.Owner.Name
				}
				
				entityList.Entities[i] = pages.EntityValues{
					ID:     note.ID,
					Values: []string{
						note.Title,
						content,
						ownerName,
						string(note.Visibility),
						note.CreatedAt.Format("2006-01-02 15:04"),
					},
				}
			}
			
			// Check if this is an HTMX request specifically for search functionality
			// Only return table content if it's an HTMX request with a search parameter or targeting the table container
			if ctx.Request().Header.Get("HX-Request") == "true" && 
				(ctx.QueryParam("search") != "" || ctx.Request().Header.Get("HX-Target") == "#user-table-container") {
				// Return only the table content for HTMX search requests
				return pages.AdminEntityListTable(ctx, n.Name, entityList)
			}
			
			return pages.AdminEntityList(ctx, n.Name, entityList)
		default:
			return echo.NewHTTPError(http.StatusNotImplemented, "entity type not supported")
		}
	}
}

func (h *Admin) EntityView(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		entity := ctx.Get(context.AdminEntityKey).(map[string][]string)
		id := ctx.Get(context.AdminEntityIDKey).(int)
		return pages.AdminEntityView(ctx, n.Name, entity, id)
	}
}

func (h *Admin) EntityAdd(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return pages.AdminEntityInput(ctx, h.getEntitySchema(n), nil)
	}
}

func (h *Admin) EntityAddSubmit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		switch n.Name {
		case "User":
			name := ctx.FormValue("name")
			phoneNumber := ctx.FormValue("phone_number")
			email := ctx.FormValue("email")
			password := ctx.FormValue("password")
			verified := ctx.FormValue("verified") == "true"
			admin := ctx.FormValue("admin") == "true"
			
			userCreate := h.orm.User.Create().
				SetName(name).
				SetPhoneNumber(phoneNumber).
				SetVerified(verified).
				SetAdmin(admin)
			
			if email != "" {
				userCreate = userCreate.SetEmail(email)
			}
			if password != "" {
				userCreate = userCreate.SetPassword(password)
			}
			
			_, err := userCreate.Save(ctx.Request().Context())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user: "+err.Error())
			}
			
			return ctx.Redirect(http.StatusSeeOther, "/admin/entity/user")
		default:
			return echo.NewHTTPError(http.StatusNotImplemented, "entity type not supported")
		}
	}
}

func (h *Admin) EntityEdit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		v := ctx.Get(context.AdminEntityKey).(map[string][]string)
		return pages.AdminEntityInput(ctx, h.getEntitySchema(n), v)
	}
}

func (h *Admin) EntityEditSubmit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid entity ID")
		}
		
		switch n.Name {
		case "User":
			name := ctx.FormValue("name")
			phoneNumber := ctx.FormValue("phone_number")
			email := ctx.FormValue("email")
			password := ctx.FormValue("password")
			verified := ctx.FormValue("verified") == "true"
			admin := ctx.FormValue("admin") == "true"
			
			userUpdate := h.orm.User.UpdateOneID(id).
				SetName(name).
				SetPhoneNumber(phoneNumber).
				SetVerified(verified).
				SetAdmin(admin)
			
			if email != "" {
				userUpdate = userUpdate.SetEmail(email)
			} else {
				userUpdate = userUpdate.ClearEmail()
			}
			if password != "" {
				// Hash the password before storing
				hashedPassword, err := h.auth.HashPassword(password)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password: "+err.Error())
				}
				userUpdate = userUpdate.SetPassword(hashedPassword)
			}
			
			_, err := userUpdate.Save(ctx.Request().Context())
			if err != nil {
				if ent.IsNotFound(err) {
					return echo.NewHTTPError(http.StatusNotFound, "user not found")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user: "+err.Error())
			}
			
			return ctx.Redirect(http.StatusSeeOther, "/admin/entity/user")
		default:
			return echo.NewHTTPError(http.StatusNotImplemented, "entity type not supported")
		}
	}
}

func (h *Admin) EntityDelete(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return pages.AdminEntityDelete(ctx, n.Name)
	}
}

func (h *Admin) EntityDeleteSubmit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid entity ID")
		}
		
		switch n.Name {
		case "User":
			err := h.orm.User.DeleteOneID(id).Exec(ctx.Request().Context())
			if err != nil {
				if ent.IsNotFound(err) {
					return echo.NewHTTPError(http.StatusNotFound, "user not found")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete user: "+err.Error())
			}
			
			return ctx.Redirect(http.StatusSeeOther, "/admin/entity/user")
		case "Note":
			err := h.orm.Note.DeleteOneID(id).Exec(ctx.Request().Context())
			if err != nil {
				if ent.IsNotFound(err) {
					return echo.NewHTTPError(http.StatusNotFound, "note not found")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete note: "+err.Error())
			}
			
			return ctx.Redirect(http.StatusSeeOther, "/admin/entity/note")
		default:
			return echo.NewHTTPError(http.StatusNotImplemented, "entity type not supported")
		}
	}
}

func (h *Admin) getEntitySchema(n *gen.Type) *load.Schema {
	for _, s := range h.graph.Schemas {
		if s.Name == n.Name {
			return s
		}
	}
	return nil
}

// VerifyUser handles POST /admin/user/:id/verify to manually verify a user account
func (h *Admin) VerifyUser(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	// Update user as verified and clear verification code
	_, err = h.orm.User.UpdateOneID(id).
		SetVerified(true).
		ClearVerificationCode().
		Save(ctx.Request().Context())
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to verify user")
	}

	return ctx.JSON(http.StatusOK, map[string]string{"status": "success", "message": "User verified successfully"})
}

func (h *Admin) Backlite(handler func(http.ResponseWriter, *http.Request) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		if id := c.Param("id"); id != "" {
			c.Request().SetPathValue("task", id)
		}
		return handler(c.Response().Writer, c.Request())
	}
}
