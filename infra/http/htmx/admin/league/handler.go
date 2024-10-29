package league

import (
	"liga-app/infra/http/htmx/common"
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

	vmodel := VModelPageLeagueDetail{
		User:   user,
		League: league,
	}

	page := PageLeagueDetail(vmodel)
	h.RenderView(c, page)
}
