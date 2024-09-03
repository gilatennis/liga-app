package team

import (
	"liga-app/infra/http/htmx/common"
	iService "liga-app/interfaces/services"
	"net/http"

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
		VModelFormTeam: &VModelFormTeam{
			Name: team.Name,
		},
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
