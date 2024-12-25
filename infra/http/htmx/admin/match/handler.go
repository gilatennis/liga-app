package match

import (
	"fmt"
	"liga-app/entities"
	"liga-app/enums"
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

	p := iCommon.NewPagination(1, 10, []iCommon.OrderBy{
		{Column: "rounds.number", Direction: "asc"},
		{Column: "matches.uuid", Direction: "asc"},
	})
	f := &services.MatchFilter{LeagueUUID: &leagueUUID, GroupUUID: &groupUUID}

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
	p := iCommon.NewPagination(page, 10, []iCommon.OrderBy{
		{Column: "rounds.number", Direction: "asc"},
		{Column: "matches.uuid", Direction: "asc"},
	})
	f := &services.MatchFilter{LeagueUUID: &leagueUUID, GroupUUID: &groupUUID}

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

func (h *Handler) PostMatchSet(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	matchUUID := params[0]

	type MatchSetPayload struct {
		PartaiName     string `form:"partai_name"`
		MatchSetNumber int    `form:"match_set_number"`
		TeamAScore     int    `form:"team_a_score"`
		TeamBScore     int    `form:"team_b_score"`
	}
	p := MatchSetPayload{}
	err = c.Bind(&p)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	matchSet := entities.MatchSet{
		MatchUUID:  matchUUID,
		Number:     p.MatchSetNumber,
		TeamAScore: p.TeamAScore,
		TeamBScore: p.TeamBScore,
	}

	_, err = h.matchService.SaveMatchSet(c, user.UUID, &matchSet)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("p.PartaiName", p.PartaiName)

	view := ViewMatchSet(user, p.PartaiName, match.LeagueUUID, match, p.MatchSetNumber, true)

	h.ResponseAlertSuccess(c, "Skor berhasil di update")
	h.RenderView(c, view)
}

func (h *Handler) PutMatchComplete(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID := params[0], params[1]
	_, _ = leagueUUID, matchUUID

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.matchService.CompleteMatch(c, user.UUID, matchUUID)
	if err != nil {
		h.ResponseAlertError(c, err.Error())
		view := ViewMatchDetailIndividuals(user, league, match)
		h.RenderView(c, view)
		return
	}

	// Fetch updated match entity
	match, err = h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	if league.ParticipationCategory == enums.CATEGORY_PARTICIPATION_TEAMS {
		view := ViewMatchDetailTeams(user, league, match)
		h.ResponseAlertSuccess(c, "Pertandingan berhasil di update")
		h.RenderView(c, view)
		return
	}

	view := ViewMatchDetailIndividuals(user, league, match)
	h.ResponseAlertSuccess(c, "Pertandingan berhasil di update")
	h.RenderView(c, view)
}

type MatchStatusPayload struct {
	MatchScheduleStatus string   `form:"match_schedule_status" json:"match_schedule_status"`
	MatchCompleteStatus []string `form:"match_complete_status[]" json:"match_complete_status"`
	MatchStatus         string   `form:"match_status" json:"match_status"`
}

func (m *MatchStatusPayload) IsMatchScheduleConfirmed() bool {
	return m.MatchScheduleStatus == "is_schedule_confirmed"
}

func (m *MatchStatusPayload) IsMatchCompleted() bool {
	for _, s := range m.MatchCompleteStatus {
		if s == "is_match_completed" {
			return true
		}
	}
	return false
}

func (m *MatchStatusPayload) IsMatchCompleteConfirmed() bool {
	for _, s := range m.MatchCompleteStatus {
		if s == "is_match_complete_confirmed" {
			return true
		}
	}
	return false
}

func (h *Handler) PutMatchStatuses(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	matchUUID := params[0]

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	p := MatchStatusPayload{}
	err = c.Bind(&p)
	if err != nil {
		view := ViewFormMatchStatuses(match)
		h.ResponseAlertError(c, err.Error())
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	isMatchScheduleConfirmed := p.IsMatchScheduleConfirmed()
	isMatchCompleted := p.IsMatchCompleted()
	isMatchCompleteConfirmed := p.IsMatchCompleteConfirmed()

	if match.MatchScheduleRequest != nil {
		match.MatchScheduleRequest.IsAccepted = isMatchScheduleConfirmed
	}
	match.IsComplete = isMatchCompleted
	match.IsCompleteConfirmed = isMatchCompleteConfirmed

	if p.MatchStatus != "" {
		match.Status = p.MatchStatus
	}

	_, err = h.matchService.Update(c, match)
	if err != nil {
		view := ViewFormMatchStatuses(match)
		h.ResponseAlertError(c, err.Error())
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	fmt.Println("isAccepted", match.MatchScheduleRequest.IsAccepted)

	h.ResponseAlertSuccess(c, "update match statuses success")
	view := ViewFormMatchStatuses(match)
	h.RenderView(c, view)
}
