package invite

import (
	"errors"
	"fmt"
	"liga-app/infra/http/htmx/common"
	"liga-app/interfaces/services"
	"net/http"

	domainErrors "liga-app/domain/errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	common.BaseHandler
	inviteService services.IInviteService
	leagueService services.ILeagueService
	teamService   services.ITeamService
}

func NewHandler(
	inviteService services.IInviteService,
	leagueService services.ILeagueService,
	teamService services.ITeamService,
) *Handler {
	return &Handler{
		inviteService: inviteService,
		leagueService: leagueService,
		teamService:   teamService,
	}
}

func (h *Handler) HandleInviteLink(c *gin.Context) {
	isLoggedIn := true

	session := sessions.Default(c)
	userUUIDAny := session.Get("user_uuid")
	userUUIDStr, ok := session.Get("user_uuid").(string)

	if userUUIDAny == nil || !ok {
		isLoggedIn = false
	}

	_ = userUUIDStr

	if !isLoggedIn {
		c.Redirect(http.StatusFound, "/auth/login")
		return

		// TODO: handle not logged in flow - refer to https://app.clickup.com/t/86eq5tnrj
		// maybe just store the redirect url client side on local storage
	}

	params, err := h.GetQueryUUIDParams(c, []string{"league_uuid", "team_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, teamUUID := params[0], params[1]

	c.Redirect(http.StatusFound, fmt.Sprintf("/invite/accept?league_uuid=%s&team_uuid=%s", leagueUUID, teamUUID))
}

func (h *Handler) GetAcceptInvitePage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetQueryUUIDParams(c, []string{"league_uuid", "team_uuid"})
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

	vmodel := &VModelPageAcceptInvite{
		User:   user,
		League: league,
		Team:   team,
	}
	page := PageAcceptInvite(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) PostAcceptInvite(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	type Data struct {
		LeagueUUID string `form:"league_uuid"`
		TeamUUID   string `form:"team_uuid"`
	}

	data := Data{}

	err = c.Bind(&data)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, err := uuid.Parse(data.LeagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	teamUUID, err := uuid.Parse(data.TeamUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.inviteService.AcceptLeagueTeamInvite(c, user.UUID, leagueUUID, teamUUID)
	var errConflict *domainErrors.ErrConflict
	if err != nil && errors.As(err, &errConflict) {
		h.ResponseAlertError(c, err.Error())
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to accept invite link")
		return
	}

	location := fmt.Sprintf("/dashboard/leagues/%s/teams/%s", leagueUUID, teamUUID)
	h.ResponseSuccessRedirect(c, "Berhasil menerima undangan tim", location)
}
