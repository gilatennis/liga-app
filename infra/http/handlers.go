package http

import (
	"fmt"
	"liga-app/config"
	"liga-app/infra/http/htmx/auth"
	"liga-app/infra/http/htmx/dashboard"
	dashLeague "liga-app/infra/http/htmx/dashboard/league"
	dashLeagueTeam "liga-app/infra/http/htmx/dashboard/league/team"
	"liga-app/infra/http/htmx/index"
	"liga-app/infra/http/htmx/invite"
	"liga-app/infra/http/json/duitku"
	"liga-app/infra/minio"
	"liga-app/infra/mysql"
	"liga-app/third_party/mailtrap"

	domainAuth "liga-app/domain/auth"
	domainInvite "liga-app/domain/invite"
	domainLeague "liga-app/domain/league"
	domainRuleBook "liga-app/domain/rule_book"
	domainTeam "liga-app/domain/team"
	domainUserInfo "liga-app/domain/user_info"
	dbAuth "liga-app/infra/mysql/auth"
	dbLeague "liga-app/infra/mysql/league"
	dbLeagueReg "liga-app/infra/mysql/league_registration"
	dbTeam "liga-app/infra/mysql/team"
	dbUser "liga-app/infra/mysql/user"
	dbWilayah "liga-app/infra/mysql/wilayah"

	thirdPartyDuitku "liga-app/third_party/duitku"

	// "liga-app/mixins"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Handlers struct {
	HTMXHandlers *HTMXHandlers
	JSONHandlers *JSONHandlers
}

type HTMXHandlers struct {
	IndexHandler               *index.Handler
	AuthHandler                *auth.Handler
	InviteHandler              *invite.Handler
	DashboardHandler           *dashboard.Handler
	DashboardLeagueHandler     *dashLeague.Handler
	DashboardLeagueTeamHandler *dashLeagueTeam.Handler
}

type JSONHandlers struct {
	DuitkuHandler *duitku.Handler
}

func setupHandlers(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger) (*Handlers, error) {
	mail := mailtrap.NewMailTrap(cfg)
	minioClient, err := minio.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("http.server.setupHandlers(): failed to create new minio client: %w", err)
	}
	_ = mail
	_ = minioClient

	duitkuAPI := thirdPartyDuitku.NewAPI(cfg)

	// fileUploader := mixins.NewFileUploader(minioClient)

	/*
		accessControlProvider, err := casbin.NewProvider(db)
		if err != nil {
			return nil, fmt.Errorf("casbin.NewProvider(db) %w", err)
		}
		accessControlStore := accessControlDB.NewStore(db)
	*/

	gormStore := mysql.NewGormStore(db)
	authStore := dbAuth.NewStore(db)
	userStore := dbUser.NewStore(db)
	wilayahStore := dbWilayah.NewStore(db)
	leagueStore := dbLeague.NewStore(db)
	leagueRegStore := dbLeagueReg.NewStore(db)
	teamStore := dbTeam.NewStore(db)

	authService := domainAuth.NewService(
		cfg,
		authStore,
		userStore,
		mail,
	)
	userInfoService := domainUserInfo.NewService(wilayahStore, userStore)
	leagueService := domainLeague.NewService(
		cfg,
		gormStore,
		leagueStore,
		leagueRegStore,
		wilayahStore,
		duitkuAPI,
	)
	teamService := domainTeam.NewService(cfg, teamStore)
	inviteService := domainInvite.NewService(leagueRegStore, leagueStore, teamStore)
	ruleBookService := domainRuleBook.NewService(gormStore)

	HTMXHandlers := HTMXHandlers{
		AuthHandler:                auth.NewHandler(cfg, authService),
		IndexHandler:               index.NewHandler(ruleBookService),
		InviteHandler:              invite.NewHandler(inviteService, leagueService, teamService),
		DashboardHandler:           dashboard.NewHandler(userInfoService, leagueService),
		DashboardLeagueHandler:     dashLeague.NewHandler(leagueService, teamService),
		DashboardLeagueTeamHandler: dashLeagueTeam.NewHandler(leagueService, teamService),
	}

	JSONHandlers := JSONHandlers{
		DuitkuHandler: duitku.NewHandler(leagueService),
	}

	return &Handlers{
		HTMXHandlers: &HTMXHandlers,
		JSONHandlers: &JSONHandlers,
	}, nil
}
