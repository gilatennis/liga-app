package fixture

import (
	"fmt"
	"liga-app/entities"
	"liga-app/infra/http/htmx/common"
	iCommon "liga-app/interfaces/common"
	"liga-app/interfaces/repositories"
	"liga-app/interfaces/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	common.BaseHandler
	fixtureService services.IFixtureService
	leagueService  services.ILeagueService
	groupService   services.IGroupService
	teamService    services.ITeamService
	matchService   services.IMatchService
}

func NewHandler(
	fixtureService services.IFixtureService,
	leagueService services.ILeagueService,
	groupService services.IGroupService,
	teamService services.ITeamService,
	matchService services.IMatchService,
) *Handler {
	return &Handler{
		fixtureService: fixtureService,
		leagueService:  leagueService,
		groupService:   groupService,
		teamService:    teamService,
		matchService:   matchService,
	}
}

func (h *Handler) GetFixtureListPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID := params[0]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	p := iCommon.NewPagination(1, 10, []iCommon.OrderBy{})
	f := repositories.FixtureFilter{LeagueUUID: &leagueUUID}
	fixtures, _, err := h.fixtureService.Get(c, p, &f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageFixtureList{
		User:     user,
		League:   league,
		Fixtures: fixtures,
	}
	page := PageFixtureList(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetFixtureDetailPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID := params[0], params[1]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	fixture, err := h.fixtureService.GetByUUID(c, fixtureUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageFixtureDetail{
		User:    user,
		League:  league,
		Fixture: fixture,
	}

	page := PageFixtureDetail(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetFixtureAddGroupPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID := params[0], params[1]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	fixture, err := h.fixtureService.GetByUUID(c, fixtureUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageAddGroup{
		User:               user,
		League:             league,
		Fixture:            fixture,
		VModelFormAddGroup: &VModelFormAddGroup{},
	}

	page := PageAddGroup(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) PostAddGroup(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID := params[0], params[1]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	fixture, err := h.fixtureService.GetByUUID(c, fixtureUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageAddGroup{
		User:               user,
		League:             league,
		Fixture:            fixture,
		VModelFormAddGroup: &VModelFormAddGroup{},
	}

	err = c.Bind(vmodel.VModelFormAddGroup)
	if err != nil {
		view := ViewFormAddGroup(vmodel)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	if !vmodel.VModelFormAddGroup.Validate() {
		view := ViewFormAddGroup(vmodel)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	group := entities.Group{
		FixtureUUID: fixtureUUID,
		Name:        vmodel.VModelFormAddGroup.Name,
	}

	_, err = h.groupService.Create(c, &group)
	if err != nil {
		vmodel.VModelFormAddGroup.base.ServerError = err.Error()
		view := ViewFormAddGroup(vmodel)
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	loc := fmt.Sprintf("/admin/dashboard/leagues/%s/fixtures/%s", leagueUUID, fixtureUUID)
	h.ResponseSuccessRedirect(c, "add group success", loc)
	view := ViewFormAddGroup(vmodel)
	h.RenderView(c, view)
}

func (h *Handler) GetGroupDetailPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid", "group_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID, groupUUID := params[0], params[1], params[2]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	fixture, err := h.fixtureService.GetByUUID(c, fixtureUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	group, err := h.groupService.GetByUUID(c, groupUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	p := iCommon.NewPagination(1, 500, []iCommon.OrderBy{{Column: "name", Direction: "asc"}})
	notingroup := true
	f := &repositories.TeamFilter{LeagueUUID: &leagueUUID, NotInGroup: &notingroup}
	teamsAvailable, _, err := h.teamService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageGroupDetail{
		User:              user,
		League:            league,
		Fixture:           fixture,
		Group:             group,
		TeamsAvailable:    teamsAvailable,
		VModelFormAddTeam: &VModelFormAddTeam{},
	}
	page := PageGroupDetail(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) PostAddTeamToGroup(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid", "group_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID, groupUUID := params[0], params[1], params[2]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	fixture, err := h.fixtureService.GetByUUID(c, fixtureUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	group, err := h.groupService.GetByUUID(c, groupUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	p := iCommon.NewPagination(1, 500, []iCommon.OrderBy{{Column: "name", Direction: "asc"}})
	notingroup := true
	f := &repositories.TeamFilter{LeagueUUID: &leagueUUID, NotInGroup: &notingroup}
	teamsAvailable, _, err := h.teamService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageGroupDetail{
		User:              user,
		League:            league,
		Fixture:           fixture,
		Group:             group,
		TeamsAvailable:    teamsAvailable,
		VModelFormAddTeam: &VModelFormAddTeam{},
	}

	err = c.Bind(vmodel.VModelFormAddTeam)
	if err != nil {
		view := ViewFormAddTeamToGroup(vmodel)
		h.ResponseAlertError(c, err.Error())
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	teamUUID, err := uuid.Parse(vmodel.VModelFormAddTeam.TeamUUID)
	if err != nil {
		view := ViewFormAddTeamToGroup(vmodel)
		h.ResponseAlertError(c, err.Error())
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	_, err = h.groupService.AddTeam(c, groupUUID, teamUUID)
	if err != nil {
		view := ViewFormAddTeamToGroup(vmodel)
		h.ResponseAlertError(c, err.Error())
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	p = iCommon.NewPagination(1, 500, []iCommon.OrderBy{{Column: "name", Direction: "asc"}})
	notingroup = true
	f = &repositories.TeamFilter{LeagueUUID: &leagueUUID, NotInGroup: &notingroup}
	vmodel.TeamsAvailable, _, err = h.teamService.Get(c, p, f)
	if err != nil {
		view := ViewFormAddTeamToGroup(vmodel)
		h.ResponseAlertError(c, err.Error())
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	vmodel.Group, err = h.groupService.GetByUUID(c, groupUUID)
	if err != nil {
		view := ViewFormAddTeamToGroup(vmodel)
		h.ResponseAlertError(c, err.Error())
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	view := ViewFormAddTeamToGroup(vmodel)
	h.ResponseAlertSuccess(c, "add team to group success")
	h.RenderView(c, view)
}

func (h *Handler) DeleteRemoveTeamFromGroup(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid", "group_uuid", "team_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID, groupUUID, teamUUID := params[0], params[1], params[2], params[3]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	fixture, err := h.fixtureService.GetByUUID(c, fixtureUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	group, err := h.groupService.GetByUUID(c, groupUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	p := iCommon.NewPagination(1, 500, []iCommon.OrderBy{{Column: "name", Direction: "asc"}})
	notingroup := true
	f := &repositories.TeamFilter{LeagueUUID: &leagueUUID, NotInGroup: &notingroup}
	teamsAvailable, _, err := h.teamService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageGroupDetail{
		User:              user,
		League:            league,
		Fixture:           fixture,
		Group:             group,
		TeamsAvailable:    teamsAvailable,
		VModelFormAddTeam: &VModelFormAddTeam{},
	}

	_, err = h.groupService.RemoveTeam(c, groupUUID, teamUUID)
	if err != nil {
		view := ViewFormAddTeamToGroup(vmodel)
		h.ResponseAlertError(c, err.Error())
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	p = iCommon.NewPagination(1, 500, []iCommon.OrderBy{{Column: "name", Direction: "asc"}})
	notingroup = true
	f = &repositories.TeamFilter{LeagueUUID: &leagueUUID, NotInGroup: &notingroup}
	vmodel.TeamsAvailable, _, err = h.teamService.Get(c, p, f)
	if err != nil {
		view := ViewFormAddTeamToGroup(vmodel)
		h.ResponseAlertError(c, err.Error())
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	vmodel.Group, err = h.groupService.GetByUUID(c, groupUUID)
	if err != nil {
		view := ViewFormAddTeamToGroup(vmodel)
		h.ResponseAlertError(c, err.Error())
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	view := ViewFormAddTeamToGroup(vmodel)
	h.ResponseAlertSuccess(c, "remove team from group success")
	h.RenderView(c, view)
}

func (h *Handler) GenerateMatches(c *gin.Context) {
	_, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid", "group_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID, groupUUID := params[0], params[1], params[2]
	_, _, _ = leagueUUID, fixtureUUID, groupUUID

	isHomeAndAwayStr := c.Query("home_and_away")
	isHomeAndAway := false

	if isHomeAndAwayStr == "true" {
		isHomeAndAway = true
	}

	err = h.matchService.GenerateRoundRobinMatches(c, groupUUID, isHomeAndAway)
	if err != nil {
		h.ResponseAlertError(c, err.Error())
		return
	}

	h.ResponseAlertSuccess(c, "round robin matches generated")
}

func (h *Handler) GenerateGroupResults(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid", "group_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID, groupUUID := params[0], params[1], params[2]
	_, _, _ = leagueUUID, fixtureUUID, groupUUID

	err = h.groupService.GenerateGroupResults(c, groupUUID)
	if err != nil {
		h.ResponseAlertError(c, err.Error())
		return
	}

	h.ResponseAlertSuccess(c, "hasil klasemen berhasil di update")
}

/*
func (h *Handler) GetGroupMatchesPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid", "group_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID, groupUUID := params[0], params[1], params[2]
	_, _, _ = leagueUUID, fixtureUUID, groupUUID

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	fixture, err := h.fixtureService.GetByUUID(c, fixtureUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	group, err := h.groupService.GetByUUID(c, groupUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	p := iCommon.NewPagination(1, 10, []iCommon.OrderBy{{Column: "rounds.number", Direction: "asc"}})
	f := &services.MatchFilter{GroupUUID: &groupUUID}

	matchStatus := c.Query("match_status")
	// matchParticipant := c.Query("match_participant")

	if matchStatus != "" && matchStatus != "ALL" {
		status := services.MATCH_FILTER_STATUS(matchStatus)
		f.Status = &status
	}

		// if matchParticipant == "" || matchParticipant == "ME" {
		// 	f.UserUUID = &user.UUID
		// }

	matches, count, err := h.matchService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}
	p.SetTotalPageByTotalItems(int(count))
	// pp.Println(matches)
	fmt.Println("p.Limit", p.Limit)

	vmodel := VModelPageGroupMatches{
		User:            user,
		League:          league,
		Fixture:         fixture,
		Group:           group,
		Matches:         matches,
		MatchesNextPage: p.GetNextPage(),
	}
	page := PageGroupMatches(vmodel)

	h.RenderView(c, page)
}

func (h *Handler) GetGroupMatchesView(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid", "group_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID, groupUUID := params[0], params[1], params[2]
	_, _, _ = leagueUUID, fixtureUUID, groupUUID

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	fixture, err := h.fixtureService.GetByUUID(c, fixtureUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	group, err := h.groupService.GetByUUID(c, groupUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	page := 1
	pageStr := c.Query("page")
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	p := iCommon.NewPagination(page, 10, []iCommon.OrderBy{{Column: "rounds.number", Direction: "asc"}})
	f := &services.MatchFilter{GroupUUID: &groupUUID}

	matchStatus := c.Query("match_status")
	// matchParticipant := c.Query("match_participant")

	if matchStatus != "" && matchStatus != "ALL" {
		status := services.MATCH_FILTER_STATUS(matchStatus)
		f.Status = &status
	}

	//
	// if matchParticipant == "ME" {
	// 	f.UserUUID = &user.UUID
	// }
	//

	matches, count, err := h.matchService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}
	p.SetTotalPageByTotalItems(int(count))

	fmt.Println("page", page)
	fmt.Println("count", count)
	vmodel := VModelPageGroupMatches{
		User:            user,
		League:          league,
		Fixture:         fixture,
		Group:           group,
		Matches:         matches,
		MatchesNextPage: p.GetNextPage(),
	}

	if page == 1 {
		view := ViewMatchList(vmodel)
		h.RenderView(c, view)
		return
	}

	if page > 1 {
		view := ViewMoreMatches(vmodel)
		h.RenderView(c, view)
		return
	}
}
*/
