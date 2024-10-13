package http

import (
	"liga-app/infra/http/middlewares"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine, handlers *Handlers, middlewares *middlewares.Middlewares) {
	// ------------------
	// static files
	// ------------------
	r.Static("/static", "./static")
	// r.LoadHTMLGlob("static/html/*")
	// r.GET("/privacy-policy", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "privacy-policy.html", gin.H{})
	// })
	// r.GET("/terms-and-conditions", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "terms-and-conditions.html", gin.H{})
	// })

	setupHtxmRoutes(r, handlers, middlewares)
	setupJsonRoutes(r, handlers, middlewares)
}

func setupHtxmRoutes(r *gin.Engine, handlers *Handlers, middlewares *middlewares.Middlewares) {
	// ------------------
	// routes serving htmx
	// ------------------
	r.GET("/", handlers.HTMXHandlers.IndexHandler.Index)

	public := r.Group("/")

	auth := public.Group("/auth")
	{
		auth.GET("/register", handlers.HTMXHandlers.AuthHandler.GetRegisterForm)
		auth.POST("/register", handlers.HTMXHandlers.AuthHandler.SubmitRegisterForm)
		auth.GET("/register/magic-link", handlers.HTMXHandlers.AuthHandler.CompleteMagicLinkRegistration)
		auth.GET("/login", handlers.HTMXHandlers.AuthHandler.GetLoginForm)
		auth.POST("/login", handlers.HTMXHandlers.AuthHandler.SubmitLoginForm)
		auth.GET("/login/magic-link", handlers.HTMXHandlers.AuthHandler.LoginWithMagicLink)
		auth.GET("/login/oauth/:provider", handlers.HTMXHandlers.AuthHandler.LoginWithOAuth)
		auth.GET("/login/oauth/:provider/callback", handlers.HTMXHandlers.AuthHandler.LoginWithOAuthCallback)
		auth.GET("/logout", handlers.HTMXHandlers.AuthHandler.LogOut)
	}

	invite := public.Group("/invite")
	{
		invite.GET("", handlers.HTMXHandlers.InviteHandler.HandleInviteLink)
		invite.GET("/accept", middlewares.SessionMiddleware, handlers.HTMXHandlers.InviteHandler.GetAcceptInvitePage)
		invite.POST("/accept", middlewares.SessionMiddleware, handlers.HTMXHandlers.InviteHandler.PostAcceptInvite)
	}

	dashboard := r.Group("/dashboard")
	dashboard.Use(middlewares.SessionMiddleware)
	{
		dashboard.GET("", handlers.HTMXHandlers.DashboardHandler.Index)
		dashboard.GET("/complete-user-info", handlers.HTMXHandlers.DashboardHandler.GetCompleteUserInfoPage)
		dashboard.POST("/complete-user-info", handlers.HTMXHandlers.DashboardHandler.PostSubmitUserInfoForm)

		leagues := dashboard.Group("/leagues")
		leagues.GET("", handlers.HTMXHandlers.DashboardLeagueHandler.GetLeagueListPage)
		leagues.GET("/:league_uuid", handlers.HTMXHandlers.DashboardLeagueHandler.GetLeagueDetailPage)

		// league registration
		leagues.GET("/:league_uuid/registration", handlers.HTMXHandlers.DashboardLeagueHandler.GetLeagueRegistrationPage)
		leagues.POST("/:league_uuid/registration", handlers.HTMXHandlers.DashboardLeagueHandler.PostLeagueRegistration)
		leagues.GET("/:league_uuid/registration/success", handlers.HTMXHandlers.DashboardLeagueHandler.GetLeagueRegistrationSuccessPage)

		// league team
		leagues.GET("/:league_uuid/teams", handlers.HTMXHandlers.DashboardLeagueTeamHandler.GetTeamListPage)
		leagues.GET("/:league_uuid/teams/view", handlers.HTMXHandlers.DashboardLeagueTeamHandler.GetTeamListView)
		leagues.GET("/:league_uuid/teams/:team_uuid", handlers.HTMXHandlers.DashboardLeagueTeamHandler.GetTeamDetailPage)
		leagues.GET("/:league_uuid/teams/:team_uuid/invites/send", handlers.HTMXHandlers.DashboardLeagueTeamHandler.GetSendInvitePage)
		leagues.GET("/:league_uuid/teams/:team_uuid/edit", handlers.HTMXHandlers.DashboardLeagueTeamHandler.GetTeamEditPage)
		leagues.PUT("/:league_uuid/teams/:team_uuid/edit", handlers.HTMXHandlers.DashboardLeagueTeamHandler.PutUpdateTeam)

		// league matches
		leagues.GET("/:league_uuid/matches", handlers.HTMXHandlers.DashboardLeagueMatchHandler.GetMatchListPage)
		leagues.GET("/:league_uuid/matches/view", handlers.HTMXHandlers.DashboardLeagueMatchHandler.GetMatchListView)

		dashboard.GET("/profile", handlers.HTMXHandlers.DashboardProfileHandler.GetProfileDetailPage)
		dashboard.GET("/profile/edit", handlers.HTMXHandlers.DashboardProfileHandler.GetProfileEditPage)
		dashboard.PUT("/profile/edit", handlers.HTMXHandlers.DashboardProfileHandler.PutUpdateProfile)
	}
}

/*
func setupCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("validOrderBy", httpCommon.ValidOrderBy)
	}
}
*/

func setupJsonRoutes(r *gin.Engine, handlers *Handlers, middlewares *middlewares.Middlewares) {
	// ------------------
	// routes serving json
	// ------------------
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// setupCustomValidators()

	apiDuitku := v1.Group("/duitku")
	{
		apiDuitku.POST("/callback", handlers.JSONHandlers.DuitkuHandler.Callback)
	}
}
