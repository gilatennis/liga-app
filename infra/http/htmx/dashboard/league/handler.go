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

	p := &iCommon.Pagination{Page: 1, Limit: 6}
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
