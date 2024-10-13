package match

import (
	"fmt"
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

	p := iCommon.NewPagination(1, 10, []iCommon.OrderBy{{Column: "rounds.number", Direction: "asc"}})
	f := &services.MatchFilter{}

	matchStatus := c.Query("match_status")
	matchParticipant := c.Query("match_participant")

	if matchStatus != "" && matchStatus != "ALL" {
		status := services.MATCH_FILTER_STATUS(matchStatus)
		f.Status = &status
	}

	if matchParticipant == "" || matchParticipant == "ME" {
		f.UserUUID = &user.UUID
	}

	matches, count, err := h.matchService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}
	p.SetTotalPageByTotalItems(int(count))
	// pp.Println(matches)
	fmt.Println("p.Limit", p.Limit)

	vmodel := VModelPageMatchList{
		User:            user,
		League:          league,
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

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID := params[0]

	page := 1
	pageStr := c.Query("page")
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	p := iCommon.NewPagination(page, 10, []iCommon.OrderBy{{Column: "rounds.number", Direction: "asc"}})
	f := &services.MatchFilter{}

	matchStatus := c.Query("match_status")
	matchParticipant := c.Query("match_participant")

	if matchStatus != "" && matchStatus != "ALL" {
		status := services.MATCH_FILTER_STATUS(matchStatus)
		f.Status = &status
	}

	if matchParticipant == "ME" {
		f.UserUUID = &user.UUID
	}

	matches, count, err := h.matchService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}
	p.SetTotalPageByTotalItems(int(count))

	fmt.Println("page", page)
	fmt.Println("count", count)

	if page == 1 {
		view := ViewMatchList(leagueUUID, matches, p.GetNextPage())
		h.RenderView(c, view)
		return
	}

	if page > 1 {
		view := ViewMoreMatches(leagueUUID, matches, p.GetNextPage())
		h.RenderView(c, view)
		return
	}
}
