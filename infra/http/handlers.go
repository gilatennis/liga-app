package http

import (
	"liga-app/config"
	admin "liga-app/infra/http/htmx/admin"
	adminAuth "liga-app/infra/http/htmx/admin/auth"
	adminCourt "liga-app/infra/http/htmx/admin/court"
	adminDash "liga-app/infra/http/htmx/admin/dashboard"
	adminLeagueFixture "liga-app/infra/http/htmx/admin/fixture"
	adminLeague "liga-app/infra/http/htmx/admin/league"
	adminLeagueTeam "liga-app/infra/http/htmx/admin/league/team"
	adminMatch "liga-app/infra/http/htmx/admin/match"
	adminPlayer "liga-app/infra/http/htmx/admin/player"
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
	domainFixture "liga-app/domain/fixture"
	domainGroup "liga-app/domain/group"
	domainInvite "liga-app/domain/invite"
	domainLeague "liga-app/domain/league"
	domainMatch "liga-app/domain/match"
	domainRuleBook "liga-app/domain/rule_book"
	domainTeam "liga-app/domain/team"
	domainUser "liga-app/domain/user"
	domainUserInfo "liga-app/domain/user_info"
	domainVenue "liga-app/domain/venue"
	dbAdmin "liga-app/infra/mysql/admin"
	dbAuth "liga-app/infra/mysql/auth"
	dbCourt "liga-app/infra/mysql/court"
	dbCourtTimeSlot "liga-app/infra/mysql/court_time_slot"
	dbFixture "liga-app/infra/mysql/fixture"
	dbGroup "liga-app/infra/mysql/group"
	dbLeague "liga-app/infra/mysql/league"
	dbLeagueReg "liga-app/infra/mysql/league_registration"
	dbMatch "liga-app/infra/mysql/match"
	dbMatchSchedReq "liga-app/infra/mysql/match_schedule_request"
	dbMatchSet "liga-app/infra/mysql/match_set"
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
	AdminHandler                *admin.Handler
	AdminDashHandler            *adminDash.Handler
	AdminAuthHandler            *adminAuth.Handler
	AdminLeagueHandler          *adminLeague.Handler
	AdminLeagueTeamHandler      *adminLeagueTeam.Handler
	AdminLeagueFixtureHandler   *adminLeagueFixture.Handler
	AdminCourtHandler           *adminCourt.Handler
	AdminPlayerHandler          *adminPlayer.Handler
	AdminMatchHandler           *adminMatch.Handler
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
	matchSetStore := dbMatchSet.NewStore(db)
	adminStore := dbAdmin.NewStore(db)
	fixtureStore := dbFixture.NewStore(db)
	groupStore := dbGroup.NewStore(db)

	authService := domainAuth.NewService(
		cfg,
		authStore,
		userStore,
		mail,
		adminStore,
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
		userStore,
	)
	teamService := domainTeam.NewService(cfg, teamStore, minioClient)
	inviteService := domainInvite.NewService(leagueRegStore, leagueStore, teamStore)
	ruleBookService := domainRuleBook.NewService(gormStore)
	matchService := domainMatch.NewService(matchStore, courtTimeSlotStore, matchSchedReqStore, matchSetStore, groupStore)
	venueService := domainVenue.NewService(venueStore)
	courtService := domainCourt.NewService(courtStore, courtTimeSlotStore)
	fixtureService := domainFixture.NewService(fixtureStore)
	groupService := domainGroup.NewService(groupStore, fixtureStore, matchStore)

	HTMXHandlers := HTMXHandlers{
		AuthHandler:                 auth.NewHandler(cfg, authService, userService),
		IndexHandler:                index.NewHandler(ruleBookService),
		InviteHandler:               invite.NewHandler(inviteService, leagueService, teamService),
		DashboardHandler:            dashboard.NewHandler(userInfoService, leagueService, matchService),
		DashboardLeagueHandler:      dashLeague.NewHandler(leagueService, teamService, matchService, groupService),
		DashboardLeagueMatchHandler: dashLeagueMatch.NewHandler(cfg, leagueService, matchService, venueService, courtService, groupService, teamService),
		DashboardLeagueTeamHandler:  dashLeagueTeam.NewHandler(leagueService, teamService),
		DashboardProfileHandler:     dashProfile.NewHandler(userInfoService),
		AdminHandler:                admin.NewHandler(),
		AdminDashHandler:            adminDash.NewHandler(leagueService),
		AdminAuthHandler:            adminAuth.NewHandler(cfg, authService, userService),
		AdminLeagueHandler:          adminLeague.NewHandler(leagueService),
		AdminLeagueTeamHandler:      adminLeagueTeam.NewHandler(leagueService, teamService, userService),
		AdminLeagueFixtureHandler:   adminLeagueFixture.NewHandler(fixtureService, leagueService, groupService, teamService, matchService),
		AdminCourtHandler:           adminCourt.NewHandler(courtService),
		AdminPlayerHandler:          adminPlayer.NewHandler(userService, userInfoService),
		AdminMatchHandler:           adminMatch.NewHandler(leagueService, matchService),
	}

	JSONHandlers := JSONHandlers{
		DuitkuHandler: duitku.NewHandler(leagueService),
	}

	return &Handlers{
		HTMXHandlers: &HTMXHandlers,
		JSONHandlers: &JSONHandlers,
	}, nil
}
