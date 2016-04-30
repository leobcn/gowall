package main

import (
	"github.com/gin-gonic/gin"
)

func BindRoutes(router *gin.Engine) {

	//front end
	router.GET("/", Index)
	router.GET("/about/", About)
	router.GET("/contact/", ContactRender)
	router.POST("/contact/", ContactSend)

	//sign up
	router.GET("/signup/", SignupRender)
	router.POST("/signup/", Signup)

	//social sign up
	router.POST("/signup/social/", SignUpSocial)
	router.GET("/signup/:provider/", SignUpProvider)

	//OAuth callback for /signup/  /login/  /settings/connect/
	router.GET("/provider/:provider/callback/", CompleteOAuth)

	//social login
	router.GET("/provider/:provider/", LoginProvider)

	//login/out
	router.GET("/login/", LoginRender)
	router.POST("/login/", Login)
	router.GET("/login/forgot/", ForgotRender)
	router.POST("/login/forgot/", SendReset)
	router.GET("/login/reset/", ResetRender)
	router.GET("/login/reset/:email/:token/", ResetRender)
	router.PUT("/login/reset/:email/:token/", ResetPassword)
	router.GET("/logout/", Logout)

	//admin
	admin := router.Group("/admin")
	admin.Use(EnsureAuthenticated)
	admin.Use(EnsureAdmin)
	{
		admin.GET("/", AdminRender)

		//admin > users
		admin.GET("/users/", AdminUsersRender)
		admin.POST("/users/", CreateUser)
		admin.GET("/users/:id/", UsersRender)
		//admin.PUT("/users/:id/", ChangeUserData)
		admin.PUT("/users/:id/password/", ChangeUserPassword)
		admin.PUT("/users/:id/role-admin/", Index)
		admin.DELETE("/users/:id/role-admin/", Index)
		admin.PUT("/users/:id/role-account/", Index)
		admin.DELETE("/users/:id/role-account/", Index)
		admin.DELETE("/users/:id/", DeleteUser)

		//admin > administrators
		admin.GET("/administrators/", renderAdministrators)
		admin.POST("/administrators/", createAdministrator)
		admin.GET("/administrators/:id/", readAdministrator)
		admin.PUT("/administrators/:id/", updateAdministrator)
		admin.PUT("/administrators/:id/permissions/", updateAdministratorPermissions)
		admin.PUT("/administrators/:id/groups/", Index)
		admin.PUT("/administrators/:id/user/", Index)
		admin.DELETE("/administrators/:id/user/", Index)
		admin.DELETE("/administrators/:id/", deleteAdministrator)

		//admin > admin groups
		admin.GET("/admin-groups/", renderAdminGroups)
		admin.POST("/admin-groups/", createAdminGroup)
		admin.GET("/admin-groups/:id/", readAdminGroup)
		admin.PUT("/admin-groups/:id/", updateAdminGroup)
		admin.PUT("/admin-groups/:id/permissions/", updateAdminGroupPermissions)
		admin.DELETE("/admin-groups/:id/", deleteAdminGroup)

		//admin > accounts
		admin.GET("/accounts/", Index)
		admin.POST("/accounts/", Index)
		admin.GET("/accounts/:id/", Index)
		admin.PUT("/accounts/:id/", Index)
		admin.PUT("/accounts/:id/user/", Index)
		admin.DELETE("/accounts/:id/user/", Index)
		admin.POST("/accounts/:id/notes/", Index)
		admin.POST("/accounts/:id/status/", Index)
		admin.DELETE("/accounts/:id/", Index)

		//admin > statuses
		admin.GET("/statuses/", renderStatuses)
		admin.POST("/statuses/", createStatus)
		admin.GET("/statuses/:id/", readStatus)
		admin.PUT("/statuses/:id/", updateStatus)
		admin.DELETE("/statuses/:id/", deleteStatus)

		//admin > categories
		admin.GET("/categories/", AdminCategoriesRender)
		admin.POST("/categories/", CreateCategory)
		admin.GET("/categories/:id/", CategoryRender)
		admin.PUT("/categories/:id/", updateCategory)
		admin.DELETE("/categories/:id/", deleteCategory)

		//admin > search
		admin.GET("/search/", Index)
	}

	//account
	account := router.Group("/account")
	account.Use(EnsureAuthenticated)
	account.Use(EnsureAccount)
	{
		account.GET("/", AccountRender)

		//account > verification
		account.GET("/verification/", AccountVerificationRender)
		account.POST("/verification/", ResendVerification)
		account.GET("/verification/:token/", Verify)

		//account > settings
		account.GET("/settings/", AccountSettingsRender)
		account.PUT("/settings/", SetSettings)
		account.PUT("/settings/identity/", ChangeIdentity)
		account.PUT("/settings/password/", ChangePassword)

		//account > settings > social
		account.GET("/providerSettings/:provider/", providerSettings)
		account.GET("/providerSettings/:provider/disconnect/", disconnectProvider)
	}

	//route not found
	router.NoRoute(Status404Render)
}