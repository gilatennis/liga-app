package http

import (
	"liga-app/config"
	"liga-app/infra/http/htmx/auth"
	"liga-app/infra/http/htmx/dashboard"
	dashLeague "liga-app/infra/http/htmx/dashboard/league"
	dashLeagueMatch "liga-app/infra/http/htmx/dashboard/league/match"
	dashLeagueTeam "liga-app/infra/http/htmx/dashboard/league/team"
	dashProfile "liga-app/infra/http/htmx/dashboard/profile"
	"liga-app/infra/http/htmx/index"
	"liga-app/infra/http/htmx/invite"
	"liga-app/infra/http/json/duitku"
	"liga-app/infra/minio"
	"liga-app/infra/mysql"
	"liga-app/third_party/mailtrap"

	domainAuth "liga-app/domain/auth"
	domainCourt "liga-app/domain/court"
	domainInvite "liga-app/domain/invite"
	domainLeague "liga-app/domain/league"
	domainMatch "liga-app/domain/match"
	domainRuleBook "liga-app/domain/rule_book"
	domainTeam "liga-app/domain/team"
	domainUser "liga-app/domain/user"
	domainUserInfo "liga-app/domain/user_info"
	domainVenue "liga-app/domain/venue"
	dbAuth "liga-app/infra/mysql/auth"
	dbCourt "liga-app/infra/mysql/court"
	dbCourtTimeSlot "liga-app/infra/mysql/court_time_slot"
	dbLeague "liga-app/infra/mysql/league"
	dbLeagueReg "liga-app/infra/mysql/league_registration"
	dbMatch "liga-app/infra/mysql/match"
	dbMatchSchedReq "liga-app/infra/mysql/match_schedule_request"
	dbTeam "liga-app/infra/mysql/team"
	dbUser "liga-app/infra/mysql/user"
	dbVenue "liga-app/infra/mysql/venue"
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
	IndexHandler                *index.Handler
	AuthHandler                 *auth.Handler
	InviteHandler               *invite.Handler
	DashboardHandler            *dashboard.Handler
	DashboardLeagueHandler      *dashLeague.Handler
	DashboardLeagueTeamHandler  *dashLeagueTeam.Handler
	DashboardLeagueMatchHandler *dashLeagueMatch.Handler
	DashboardProfileHandler     *dashProfile.Handler
}

type JSONHandlers struct {
	DuitkuHandler *duitku.Handler
}

func setupHandlers(
	cfg *config.Config,
	db *gorm.DB,
	minioClient *minio.Client,
	logger *zerolog.Logger,
) (*Handlers, error) {
	mail := mailtrap.NewMailTrap(cfg)
	_ = mail
	_ = logger

	duitkuAPI := thirdPartyDuitku.NewAPI(cfg)

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
	courtTimeSlotStore := dbCourtTimeSlot.NewStore(db)
	matchStore := dbMatch.NewStore(db)
	venueStore := dbVenue.NewStore(db)
	courtStore := dbCourt.NewStore(db)
	matchSchedReqStore := dbMatchSchedReq.NewStore(db)

	authService := domainAuth.NewService(
		cfg,
		authStore,
		userStore,
		mail,
	)
	userService := domainUser.NewService(userStore)
	userInfoService := domainUserInfo.NewService(
		cfg,
		wilayahStore,
		userStore,
		minioClient,
	)
	leagueService := domainLeague.NewService(
		cfg,
		gormStore,
		leagueStore,
		leagueRegStore,
		wilayahStore,
		duitkuAPI,
	)
	teamService := domainTeam.NewService(cfg, teamStore, minioClient)
	inviteService := domainInvite.NewService(leagueRegStore, leagueStore, teamStore)
	ruleBookService := domainRuleBook.NewService(gormStore)
	matchService := domainMatch.NewService(matchStore, courtTimeSlotStore, matchSchedReqStore)
	venueService := domainVenue.NewService(venueStore)
	courtService := domainCourt.NewService(courtStore, courtTimeSlotStore)

	HTMXHandlers := HTMXHandlers{
		AuthHandler:                 auth.NewHandler(cfg, authService, userService),
		IndexHandler:                index.NewHandler(ruleBookService),
		InviteHandler:               invite.NewHandler(inviteService, leagueService, teamService),
		DashboardHandler:            dashboard.NewHandler(userInfoService, leagueService),
		DashboardLeagueHandler:      dashLeague.NewHandler(leagueService, teamService),
		DashboardLeagueMatchHandler: dashLeagueMatch.NewHandler(leagueService, matchService, venueService, courtService),
		DashboardLeagueTeamHandler:  dashLeagueTeam.NewHandler(leagueService, teamService),
		DashboardProfileHandler:     dashProfile.NewHandler(userInfoService),
	}

	JSONHandlers := JSONHandlers{
		DuitkuHandler: duitku.NewHandler(leagueService),
	}

	return &Handlers{
		HTMXHandlers: &HTMXHandlers,
		JSONHandlers: &JSONHandlers,
	}, nil
}
