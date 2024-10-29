package team

import (
	"fmt"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/infra/http/htmx/common"
	iCommon "liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"
	iService "liga-app/interfaces/services"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

type Handler struct {
	common.BaseHandler
	leagueService iService.ILeagueService
	teamService   iService.ITeamService
	userService   iService.IUserService
}

func NewHandler(
	leagueService iService.ILeagueService,
	teamService iService.ITeamService,
	userService iService.IUserService,
) *Handler {
	return &Handler{
		leagueService: leagueService,
		teamService:   teamService,
		userService:   userService,
	}
}

func (h *Handler) GetTeamListPage(c *gin.Context) {
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
	f := &iRepo.TeamFilter{LeagueUUID: &leagueUUID}
	teams, count, err := h.teamService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}
	p.SetTotalPageByTotalItems(int(count))

	vmodel := &VModelPageTeamList{
		User:          user,
		League:        league,
		Teams:         teams,
		TeamsNextPage: p.GetNextPage(),
	}

	page := PageTeamList(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetTeamListView(c *gin.Context) {
	_, err := h.GetAuthenticatedUser(c)
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

	type Query struct {
		Page int    `form:"page"`
		View string `form:"view"`
	}

	q := Query{}

	err = c.Bind(&q)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	p := iCommon.NewPagination(q.Page, 10, []iCommon.OrderBy{})
	f := &iRepo.TeamFilter{LeagueUUID: &leagueUUID}
	teams, count, err := h.teamService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}
	p.SetTotalPageByTotalItems(int(count))

	var view templ.Component

	switch q.View {
	case "teams":
		view = ViewTeamList(leagueUUID, teams, p.GetNextPage())
	case "players":
		view = ViewPlayerList(leagueUUID, teams, p.GetNextPage())
	default:
		h.RenderErrorPage(c, http.StatusBadRequest, "view must be teams or players")
		return
	}

	h.RenderView(c, view)
}

func (h *Handler) GetTeamDetailPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "team_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, teamUUID := params[0], params[1]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	team, err := h.teamService.GetByUUID(c, teamUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := &VModelPageTeamDetail{
		User:   user,
		League: league,
		Team:   team,
	}
	page := PageTeamDetail(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetAddTeamPage(c *gin.Context) {
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

	p := iCommon.NewPagination(1, 500, []iCommon.OrderBy{{Column: "name", Direction: "asc"}})
	players, _, err := h.userService.GetUsers(c, p, nil)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageAddTeam{
		User:              user,
		League:            league,
		Players:           players,
		VModelFormAddTeam: &VModelFormAddTeam{},
	}
	page := PageAddTeam(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) PostAddTeam(c *gin.Context) {
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

	p := iCommon.NewPagination(1, 500, []iCommon.OrderBy{{Column: "name", Direction: "asc"}})
	players, _, err := h.userService.GetUsers(c, p, nil)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageAddTeam{
		User:              user,
		League:            league,
		Players:           players,
		VModelFormAddTeam: &VModelFormAddTeam{},
	}

	err = c.Bind(vmodel.VModelFormAddTeam)
	if err != nil {
		view := ViewFormAddTeam(vmodel)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	teamName := vmodel.VModelFormAddTeam.TeamName

	pp.Println(vmodel.VModelFormAddTeam)

	playerUUIDs := []uuid.UUID{}
	for _, pUUIDStr := range vmodel.VModelFormAddTeam.PlayerUUIDs {
		pUUID, err := uuid.Parse(pUUIDStr)
		if err != nil {
			vmodel.VModelFormAddTeam.base.ServerError = err.Error()
			view := ViewFormAddTeam(vmodel)
			h.RenderViewWithStatus(c, view, http.StatusBadRequest)
			return
		}
		playerUUIDs = append(playerUUIDs, pUUID)
	}

	leagueReg := entities.LeagueRegistration{
		LeagueUUID:    leagueUUID,
		UserUUID:      playerUUIDs[0],
		PaymentAmount: vmodel.VModelFormAddTeam.PaymentAmount,
		PaymentID:     vmodel.VModelFormAddTeam.PaymentID,
		PaymentStatus: enums.PAYMENT_STATUS_PAID,
	}

	_, err = h.leagueService.RegisterTeamToLeague(c, &leagueReg, teamName, playerUUIDs)
	if err != nil {
		vmodel.VModelFormAddTeam.base.ServerError = err.Error()
		view := ViewFormAddTeam(vmodel)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	h.ResponseSuccessRedirect(c, "berhasil menambah tim / pemain", fmt.Sprintf("/admin/dashboard/leagues/%s", leagueUUID))
	h.RenderView(c, ViewFormAddTeam(vmodel))
}
