package league

import (
	"liga-app/infra/http/htmx/common"
	iCommon "liga-app/interfaces/common"
	iService "liga-app/interfaces/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	common.BaseHandler
	leagueService iService.ILeagueService
}

func NewHandler(leagueService iService.ILeagueService) *Handler {
	return &Handler{
		leagueService: leagueService,
	}
}

func (h *Handler) GetLeagueListPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	p := iCommon.NewPagination(1, 8, []iCommon.OrderBy{})
	leagues, count, err := h.leagueService.Get(c, p, nil)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	p.SetTotalPageByTotalItems(int(count))

	vmodel := &VModelPageLeagueList{
		User:            user,
		Leagues:         leagues,
		LeaguesNextPage: p.GetNextPage(),
	}
	page := PageLeagueList(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetLeagueDetailPage(c *gin.Context) {
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

	vmodel := &VModelPageLeagueDetail{
		User:   user,
		League: league,
	}

	page := PageLeagueDetail(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetLeagueRegistrationPage(c *gin.Context) {
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

	vmodel := &VModelPageLeagueRegistration{
		User:   user,
		League: league,
	}

	page := PageLeagueRegistration(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) PostLeagueRegistration(c *gin.Context) {
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

	paymentURL, err := h.leagueService.Join(c, leagueUUID, user)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("HX-Redirect", paymentURL)
}

func (h *Handler) GetLeagueRegistrationSuccessPage(c *gin.Context) {
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

	vmodel := &VModelPageLeagueRegistrationSuccess{
		User:   user,
		League: league,
	}
	page := PageLeagueRegistrationSuccess(vmodel)
	h.RenderView(c, page)
}
