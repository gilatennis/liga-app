package match

import (
	"liga-app/infra/http/htmx/common"
	iCommon "liga-app/interfaces/common"
	"liga-app/interfaces/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	common.BaseHandler
	leagueService services.ILeagueService
	matchService  services.IMatchService
}

func NewHandler(
	leagueService services.ILeagueService,
	matchService services.IMatchService,
) *Handler {
	return &Handler{
		leagueService: leagueService,
		matchService:  matchService,
	}
}

func (h *Handler) GetMatchListPage(c *gin.Context) {
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

	p := iCommon.NewPagination(1, 10, []iCommon.OrderBy{{Column: "rounds.number", Direction: "asc"}})
	f := &services.MatchFilter{LeagueUUID: &leagueUUID}

	matchStatus := c.Query("match_status")
	// matchParticipant := c.Query("match_participant")

	if matchStatus != "" && matchStatus != "ALL" {
		status := services.MATCH_FILTER_STATUS(matchStatus)
		f.Status = &status
	}

	matches, count, err := h.matchService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}
	p.SetTotalPageByTotalItems(int(count))

	vmodel := VModelPageMatchList{
		User:            user,
		LeagueUUID:      leagueUUID,
		FixtureUUID:     fixtureUUID,
		GroupUUID:       groupUUID,
		Matches:         matches,
		MatchesNextPage: p.GetNextPage(),
	}
	page := PageMatchList(&vmodel)

	h.RenderView(c, page)
}

func (h *Handler) GetMatchListView(c *gin.Context) {
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

	page := 1
	pageStr := c.Query("page")
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	p := iCommon.NewPagination(page, 10, []iCommon.OrderBy{{Column: "rounds.number", Direction: "asc"}})
	f := &services.MatchFilter{LeagueUUID: &leagueUUID}

	matchStatus := c.Query("match_status")
	matchParticipant := c.Query("match_participant")
	matchScheduleStatus := c.Query("match_schedule_status")

	if matchStatus != "" && matchStatus != "ALL" {
		status := services.MATCH_FILTER_STATUS(matchStatus)
		f.Status = &status
	}

	if matchParticipant == "ME" {
		f.UserUUID = &user.UUID
	}

	if matchScheduleStatus == "is_confirmed" {
		isConfirmed := true
		f.IsScheduleConfirmed = &isConfirmed
	}

	matches, count, err := h.matchService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}
	p.SetTotalPageByTotalItems(int(count))

	if page == 1 {
		view := ViewMatchList(leagueUUID, fixtureUUID, groupUUID, matches, p.GetNextPage())
		h.RenderView(c, view)
		return
	}

	if page > 1 {
		view := ViewMoreMatches(leagueUUID, fixtureUUID, groupUUID, matches, p.GetNextPage())
		h.RenderView(c, view)
		return
	}
}

func (h *Handler) GetMatchDetailPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "fixture_uuid", "group_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, fixtureUUID, groupUUID, matchUUID := params[0], params[1], params[2], params[3]
	_, _, _, _ = leagueUUID, fixtureUUID, groupUUID, matchUUID

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageMatchDetail{
		User:   user,
		League: league,
		Match:  match,
	}

	page := PageMatchDetail(vmodel)
	h.RenderView(c, page)
}
