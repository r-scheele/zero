package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type StudyGroups struct {
	container *services.Container
}

func init() {
	Register(new(StudyGroups))
}

func (h *StudyGroups) Init(c *services.Container) error {
	h.container = c
	return nil
}

func (h *StudyGroups) Routes(g *echo.Group) {
	// Study Groups routes (require authentication and verification)
	studyGroups := g.Group("/study-groups", middleware.RequireAuthentication, middleware.RequireVerification)
	
	// List study groups
	studyGroups.GET("", h.ListStudyGroups).Name = "study-groups.list"
	
	// Create study group
	studyGroups.GET("/create", h.CreateStudyGroupPage).Name = "study-groups.create"
	studyGroups.POST("/create", h.CreateStudyGroupSubmit).Name = "study-groups.create.submit"
	
	// View study group
	studyGroups.GET("/:id", h.ViewStudyGroup).Name = "study-groups.view"
	
	// Join study group
	studyGroups.POST("/:id/join", h.JoinStudyGroup).Name = "study-groups.join"
	
	// Leave study group
	studyGroups.POST("/:id/leave", h.LeaveStudyGroup).Name = "study-groups.leave"
}

// ListStudyGroups displays all study groups
func (h *StudyGroups) ListStudyGroups(ctx echo.Context) error {
	return pages.StudyGroupsList(ctx)
}

// CreateStudyGroupPage displays the study group creation form
func (h *StudyGroups) CreateStudyGroupPage(ctx echo.Context) error {
	return pages.CreateStudyGroup(ctx)
}

// CreateStudyGroupSubmit handles study group creation
func (h *StudyGroups) CreateStudyGroupSubmit(ctx echo.Context) error {
	// TODO: Implement study group creation logic
	return ctx.String(200, "Study group creation functionality coming soon!")
}

// ViewStudyGroup displays a specific study group
func (h *StudyGroups) ViewStudyGroup(ctx echo.Context) error {
	return pages.ViewStudyGroup(ctx)
}

// JoinStudyGroup handles joining a study group
func (h *StudyGroups) JoinStudyGroup(ctx echo.Context) error {
	// TODO: Implement join study group logic
	return ctx.String(200, "Join study group functionality coming soon!")
}

// LeaveStudyGroup handles leaving a study group
func (h *StudyGroups) LeaveStudyGroup(ctx echo.Context) error {
	// TODO: Implement leave study group logic
	return ctx.String(200, "Leave study group functionality coming soon!")
}