package match

import (
	"fmt"
	"liga-app/config"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/infra/http/htmx/common"
	iCommon "liga-app/interfaces/common"
	"liga-app/interfaces/repositories"
	"liga-app/interfaces/services"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

type Handler struct {
	common.BaseHandler
	cfg           *config.Config
	leagueService services.ILeagueService
	matchService  services.IMatchService
	venueService  services.IVenueService
	courtService  services.ICourtService
}

func NewHandler(
	cfg *config.Config,
	leagueService services.ILeagueService,
	matchService services.IMatchService,
	venueService services.IVenueService,
	courtService services.ICourtService,
) *Handler {
	return &Handler{
		cfg:           cfg,
		leagueService: leagueService,
		matchService:  matchService,
		venueService:  venueService,
		courtService:  courtService,
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
	f := &services.MatchFilter{LeagueUUID: &leagueUUID}

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
	f := &services.MatchFilter{LeagueUUID: &leagueUUID}

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

func (h *Handler) GetMatchDetailPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID := params[0], params[1]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageMatchDetail{
		User:           user,
		League:         league,
		Match:          match,
		WebsiteBaseURL: h.cfg.App.BaseURL,
	}

	page := PageMatchDetail(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetMatchDetailView(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID := params[0], params[1]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageMatchDetail{
		User:           user,
		League:         league,
		Match:          match,
		WebsiteBaseURL: h.cfg.App.BaseURL,
	}

	view := ViewMatchDetail(vmodel)
	h.RenderView(c, view)
}

func (h *Handler) GetIsMatchResultUpdated(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID := params[0], params[1]

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	lastUpdateStr := c.Query("last-update")

	fmt.Println("lastUpdateStr", lastUpdateStr)
	_ = user
	_ = match

	lastUpdate, err := time.Parse(time.RFC3339, lastUpdateStr)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("lastUpdate", lastUpdate)
	fmt.Println("match.UpdatedAt", match.UpdatedAt)

	if match.UpdatedAt.Sub(lastUpdate) > (5 * time.Second) {
		c.Header("HX-Trigger", "match-update")
	}
	view := ViewMatchResultUpdatedCheck(leagueUUID, matchUUID, match.UpdatedAt)
	h.RenderView(c, view)
}

func (h *Handler) GetMatchScheduleVenueListPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID := params[0], params[1]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	p := iCommon.NewPagination(1, 10, []iCommon.OrderBy{{Column: "rounds.number", Direction: "asc"}})
	f := &repositories.VenueFilter{}
	venues, _, err := h.venueService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageMatchScheduleVenues{
		User:   user,
		League: league,
		Match:  match,
		Venues: venues,
	}

	page := PageMatchScheduleVenues(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetUseOwnVenuePage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID := params[0], params[1]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageMatchScheduleOwnVenue{
		User:                            user,
		League:                          league,
		Match:                           match,
		VModelFormMatchScheduleOwnVenue: &VModelFormMatchScheduleOwnVenue{},
	}

	page := PageMatchScheduleOwnVenue(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) PostUseOwnVenueForm(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID := params[0], params[1]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	_ = user
	_ = league
	_ = match

	vmodelPage := VModelPageMatchScheduleOwnVenue{
		User:   user,
		League: league,
		Match:  match,
	}

	p := VModelFormMatchScheduleOwnVenue{}
	err = c.Bind(&p)
	if err != nil {
		vmodelPage.VModelFormMatchScheduleOwnVenue = &p
		view := ViewOwnVenueForm(vmodelPage)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	if !p.Validate() {
		vmodelPage.VModelFormMatchScheduleOwnVenue = &p
		view := ViewOwnVenueForm(vmodelPage)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	vmodelPage.VModelFormMatchScheduleOwnVenue = &p

	req := entities.MatchScheduleRequest{
		VenueName: p.VenueName,
		Address:   p.Address,
		Notes:     p.Notes,
		Start:     &p.Start,
	}
	_, err = h.matchService.RequestMatchScheduleOwnVenue(c, user.UUID, matchUUID, &req)
	if err != nil {
		view := ViewOwnVenueForm(vmodelPage)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	location := fmt.Sprintf("/dashboard/leagues/%s/matches/%s",
		leagueUUID, matchUUID,
	)
	h.ResponseSuccessRedirect(c, "request match schedule success", location)

	view := ViewOwnVenueForm(vmodelPage)
	h.RenderView(c, view)
}

func (h *Handler) GetMatchScheduleCourtListPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid", "venue_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID, venueUUID := params[0], params[1], params[2]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	venue, err := h.venueService.GetByUUID(c, venueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageMatchScheduleCourts{
		User:   user,
		League: league,
		Match:  match,
		Venue:  venue,
	}

	page := PageMatchScheduleCourts(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetMatchScheduleCourtDaysPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid", "venue_uuid", "court_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID, venueUUID, courtUUID := params[0], params[1], params[2], params[3]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	venue, err := h.venueService.GetByUUID(c, venueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	court, err := h.courtService.GetByUUID(c, courtUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	availableDays, err := h.courtService.GetCourtAvailableDays(c, courtUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	pp.Println("availableDays", availableDays)

	vmodel := VModelPageMatchScheduleDays{
		User:          user,
		League:        league,
		Match:         match,
		Venue:         venue,
		Court:         court,
		AvailableDays: availableDays,
	}
	page := PageMatchScheduleDays(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetMatchScheduleCourtTimesPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid", "venue_uuid", "court_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID, venueUUID, courtUUID := params[0], params[1], params[2], params[3]

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	venue, err := h.venueService.GetByUUID(c, venueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	court, err := h.courtService.GetByUUID(c, courtUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	day := c.Query("day")
	fmt.Println("day", day)

	t, err := time.Parse(time.RFC3339, day)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	timeslots, err := h.courtService.GetCourtTimeSlotsByDay(c, courtUUID, t)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageMatchScheduleTimes{
		User:      user,
		League:    league,
		Match:     match,
		Venue:     venue,
		Court:     court,
		Day:       t,
		TimeSlots: timeslots,
	}
	_ = vmodel

	page := PageMatchScheduleTimes(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) PostRequestMatchSchedule(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid", "venue_uuid", "court_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID, venueUUID, courtUUID := params[0], params[1], params[2], params[3]
	_, _, _, _ = leagueUUID, matchUUID, venueUUID, courtUUID

	type Payload struct {
		TimeSlotUUIDs string `form:"timeslot_uuids"`
	}
	p := Payload{}
	err = c.Bind(&p)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	timeslotUUIDStrings := strings.Split(p.TimeSlotUUIDs, ",")

	timeslotUUIDs := []uuid.UUID{}
	for _, s := range timeslotUUIDStrings {
		tsUUID, err := uuid.Parse(s)
		if err != nil {
			h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
			return
		}
		timeslotUUIDs = append(timeslotUUIDs, tsUUID)
	}

	scheduleRequestUUID, err := h.matchService.RequestMatchSchedule(c, user.UUID, matchUUID, timeslotUUIDs)
	if err != nil {
		h.ResponseAlertError(c, err.Error())
		return
	}
	_ = scheduleRequestUUID

	location := fmt.Sprintf("/dashboard/leagues/%s/matches/%s",
		leagueUUID, matchUUID,
	)
	h.ResponseSuccessRedirect(c, "request match schedule success", location)
}

/*
func (h *Handler) GetScheduleRequestDetailPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid", "schedule_request_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID, scheduleRequestUUID := params[0], params[1], params[2]
	_, _, _ = leagueUUID, matchUUID, scheduleRequestUUID

	fmt.Println("scheduleRequestUUID", scheduleRequestUUID)

	req, err := h.matchService.GetScheduleRequestByUUID(c, scheduleRequestUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}
}
*/

func (h *Handler) PutConfirmMatchSchedule(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid", "schedule_request_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID, scheduleRequestUUID := params[0], params[1], params[2]
	_, _, _ = leagueUUID, matchUUID, scheduleRequestUUID

	_, err = h.matchService.ConfirmMatchSchedule(c, user.UUID, scheduleRequestUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	location := fmt.Sprintf("/dashboard/leagues/%s/matches/%s", leagueUUID, matchUUID)
	h.ResponseSuccessRedirect(c, "confirm match schedule success", location)
}

func (h *Handler) PostMatchSet(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID := params[0], params[1]
	_, _ = leagueUUID, matchUUID

	type MatchSetPayload struct {
		PartaiName     string `form:"partai_name"`
		MatchSetNumber int    `form:"match_set_number"`
		TeamAScore     int    `form:"team_a_score"`
		TeamBScore     int    `form:"team_b_score"`
	}
	p := MatchSetPayload{}
	err = c.Bind(&p)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	matchSet := entities.MatchSet{
		MatchUUID:  matchUUID,
		Number:     p.MatchSetNumber,
		TeamAScore: p.TeamAScore,
		TeamBScore: p.TeamBScore,
	}

	_, err = h.matchService.SaveMatchSet(c, user.UUID, &matchSet)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	view := ViewMatchSet(user, p.PartaiName, leagueUUID, match, p.MatchSetNumber, true)

	h.ResponseAlertSuccess(c, "Skor berhasil di update")
	h.RenderView(c, view)
}

func (h *Handler) PutMatchComplete(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID := params[0], params[1]
	_, _ = leagueUUID, matchUUID

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.matchService.CompleteMatch(c, user.UUID, matchUUID)
	if err != nil {
		h.ResponseAlertError(c, err.Error())
		view := ViewMatchDetailIndividuals(user, league, match)
		h.RenderView(c, view)
		return
	}

	// Fetch updated match entity
	match, err = h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	if league.ParticipationCategory == enums.CATEGORY_PARTICIPATION_TEAMS {
		view := ViewMatchDetailTeams(user, league, match)
		h.ResponseAlertSuccess(c, "Pertandingan berhasil di update")
		h.RenderView(c, view)
		return
	}

	view := ViewMatchDetailIndividuals(user, league, match)
	h.ResponseAlertSuccess(c, "Pertandingan berhasil di update")
	h.RenderView(c, view)
}

func (h *Handler) PutMatchConfirmComplete(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID := params[0], params[1]
	_, _ = leagueUUID, matchUUID

	_, err = h.matchService.ConfirmCompleteMatch(c, user.UUID, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	league, err := h.leagueService.GetByUUID(c, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	match, err := h.matchService.GetByUUID(c, matchUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	if league.ParticipationCategory == enums.CATEGORY_PARTICIPATION_TEAMS {
		view := ViewMatchDetailTeams(user, league, match)
		h.ResponseAlertSuccess(c, "Pertandingan berhasil di update")
		h.RenderView(c, view)
		return
	}

	view := ViewMatchDetailIndividuals(user, league, match)
	h.ResponseAlertSuccess(c, "Pertandingan berhasil di update")
	h.RenderView(c, view)
}

func (h *Handler) DeleteCancelMatchSchedule(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}
	_ = user

	params, err := h.GetURLUUIDParams(c, []string{"league_uuid", "match_uuid", "schedule_request_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	leagueUUID, matchUUID, scheduleRequestUUID := params[0], params[1], params[2]
	_, _, _ = leagueUUID, matchUUID, scheduleRequestUUID
}
