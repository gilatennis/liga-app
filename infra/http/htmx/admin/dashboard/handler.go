package dashboard

import (
	"liga-app/infra/http/htmx/common"
	"net/http"

	iCommon "liga-app/interfaces/common"
	iService "liga-app/interfaces/services"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
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

func (h *Handler) Index(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	p := iCommon.NewPagination(1, 8, []iCommon.OrderBy{})
	leagues, _, err := h.leagueService.Get(c, p, nil)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	pp.Println("leagues", leagues)

	vmodel := VModelPageDashboard{
		User:    user,
		Leagues: leagues,
	}
	page := PageDashboard(vmodel)
	h.RenderView(c, page)
}
