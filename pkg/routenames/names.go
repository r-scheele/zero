package routenames

import (
	"fmt"
)

const (
	Home                  = "home"
	Dashboard             = "dashboard"
	Quizzes               = "quizzes"
	Summaries             = "summaries"
	About                 = "about"
	Contact               = "contact"
	ContactSubmit         = "contact.submit"
	Login                 = "login"
	LoginSubmit           = "login.submit"
	Register              = "register"
	RegisterSubmit        = "register.submit"
	ForgotPassword        = "forgot_password"
	ForgotPasswordSubmit  = "forgot_password.submit"
	Logout                = "logout"
	Profile               = "profile"
	ProfileEdit           = "profile.edit"
	ProfileUpdate         = "profile.update"
	ProfilePicture        = "profile.picture"
	ProfileChangePassword = "profile.change_password"
	ProfileDeactivate     = "profile.deactivate"
	VerifyEmail           = "verify_email"
	VerificationNotice    = "verification_notice"
	ResendVerification    = "resend_verification"
	ResetPassword         = "reset_password"
	ResetPasswordSubmit   = "reset_password.submit"
	Search                = "search"
	Task                  = "task"
	TaskSubmit            = "task.submit"
	Cache                 = "cache"
	CacheSubmit           = "cache.submit"
	Files                 = "files"
	FilesSubmit           = "files.submit"
	AdminTasks            = "admin:tasks"
)

func AdminEntityList(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_list", entityTypeName)
}

func AdminEntityAdd(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_add", entityTypeName)
}

func AdminEntityView(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_view", entityTypeName)
}

func AdminEntityEdit(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_edit", entityTypeName)
}

func AdminEntityDelete(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_delete", entityTypeName)
}

func AdminEntityAddSubmit(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_add.submit", entityTypeName)
}

func AdminEntityEditSubmit(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_edit.submit", entityTypeName)
}

func AdminEntityDeleteSubmit(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_delete.submit", entityTypeName)
}
