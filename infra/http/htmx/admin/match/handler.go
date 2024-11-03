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
	matchService services.IMatchService
}

func NewHandler(matchService services.IMatchService) *Handler {
	return &Handler{
		matchService: matchService,
	}
}

func (h *Handler) GetMatchListPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	p := iCommon.NewPagination(1, 10, []iCommon.OrderBy{})
	f := &services.MatchFilter{}

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

	page := 1
	pageStr := c.Query("page")
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	p := iCommon.NewPagination(page, 10, []iCommon.OrderBy{{Column: "rounds.number", Direction: "asc"}})
	f := &services.MatchFilter{}

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
		view := ViewMatchList(matches, p.GetNextPage())
		h.RenderView(c, view)
		return
	}

	if page > 1 {
		view := ViewMoreMatches(matches, p.GetNextPage())
		h.RenderView(c, view)
		return
	}
}
