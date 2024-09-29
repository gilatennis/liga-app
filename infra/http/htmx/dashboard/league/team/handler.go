package team

import (
	"fmt"
	"liga-app/infra/http/htmx/common"
	iCommon "liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"
	iService "liga-app/interfaces/services"
	"mime/multipart"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	common.BaseHandler
	leagueService iService.ILeagueService
	teamService   iService.ITeamService
}

func NewHandler(
	leagueService iService.ILeagueService,
	teamService iService.ITeamService,
) *Handler {
	return &Handler{
		leagueService: leagueService,
		teamService:   teamService,
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

func (h *Handler) GetSendInvitePage(c *gin.Context) {
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

	inviteLink := h.teamService.GetLeagueTeamInviteLink(c, leagueUUID, teamUUID)

	vmodel := &VModelPageSendInvite{
		User:                     user,
		League:                   league,
		Team:                     team,
		InviteLink:               inviteLink,
		VModelFormInviteViaEmail: &VModelFormInviteViaEmail{},
	}
	page := PageSendInvite(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetTeamEditPage(c *gin.Context) {
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

	team, err := h.teamService.GetByUUID(c, teamUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := &VModelPageTeamEdit{
		User: user,
		VModelFormTeam: &VModelFormTeam{
			LeagueUUID:            leagueUUID,
			TeamUUID:              teamUUID,
			ProfilePictureCurrent: team.ProfilePictureURL,
			Name:                  team.Name,
		},
	}

	page := PageTeamEdit(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) PutUpdateTeam(c *gin.Context) {
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

	team, err := h.teamService.GetByUUID(c, teamUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	form := VModelFormTeam{
		LeagueUUID:            leagueUUID,
		TeamUUID:              teamUUID,
		ProfilePictureCurrent: team.ProfilePictureURL,
	}

	err = c.Bind(&form)
	if err != nil {
		form.base.ServerError = err.Error()
		view := ViewFormEditTeam(&form)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	if !form.Validate() {
		// pp.Println("form.base.formErrors", form.base.FormErrors)
		view := ViewFormEditTeam(&form)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	fileHeader, err := c.FormFile("profile_picture")
	if err != nil && err.Error() != "http: no such file" {
		form.base.ServerError = err.Error()
		view := ViewFormEditTeam(&form)
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	var file *multipart.File = nil

	if fileHeader != nil {
		_file, err := fileHeader.Open()
		if err != nil {
			form.base.ServerError = err.Error()
			view := ViewFormEditTeam(&form)
			h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
			return
		}
		file = &_file
	}

	team.Name = form.Name

	_, err = h.teamService.Update(c, user, team, file)
	if err != nil {
		form.base.ServerError = err.Error()
		view := ViewFormEditTeam(&form)
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	h.ResponseSuccessRedirect(c, "informasi tim berhasil di update", fmt.Sprintf("/dashboard/leagues/%s/teams/%s", leagueUUID, teamUUID))
	h.RenderView(c, ViewFormEditTeam(&form))
}
