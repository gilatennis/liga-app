package league

import (
	"fmt"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/infra/http/htmx/common"
	iCommon "liga-app/interfaces/common"
	iService "liga-app/interfaces/services"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	common.BaseHandler
	leagueService iService.ILeagueService
	teamService   iService.ITeamService
	matchService  iService.IMatchService
	groupService  iService.IGroupService
}

func NewHandler(
	leagueService iService.ILeagueService,
	teamService iService.ITeamService,
	matchService iService.IMatchService,
	groupService iService.IGroupService,
) *Handler {
	return &Handler{
		leagueService: leagueService,
		teamService:   teamService,
		matchService:  matchService,
		groupService:  groupService,
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
		ShowPlayerCount: false,
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

	team, err := h.teamService.GetUserTeamInLeague(c, user.UUID, leagueUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	isRatingRequired := false
	numPlayersHasRating := 0

	playersWithRatingTooHigh := []entities.User{}
	playersWithRatingTooLow := []entities.User{}

	if team != nil {
		for _, p := range team.Players {
			for _, r := range p.Ratings {
				if r.Season == league.Season {
					numPlayersHasRating += 1

					if league.IsNTRPRatingTooHigh(r.NTRPRating) {
						playersWithRatingTooHigh = append(playersWithRatingTooHigh, p)
					}

					if league.IsNTRPRatingTooLow(r.NTRPRating) {
						playersWithRatingTooLow = append(playersWithRatingTooLow, p)
					}
				}
			}
		}

		if len(team.Players) > numPlayersHasRating {
			isRatingRequired = true
		}
	}

	p := iCommon.NewPagination(1, 1, []iCommon.OrderBy{{Column: "rounds.number", Direction: "asc"}})
	isComplete := true
	isCompleteConfirmed := false
	f := &iService.MatchFilter{LeagueUUID: &leagueUUID, UserUUID: &user.UUID, IsComplete: &isComplete, IsCompleteConfirmed: &isCompleteConfirmed}
	matches, _, err := h.matchService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	var matchNeedsCompleteConfirmation *entities.Match
	if len(matches) > 0 {
		matchNeedsCompleteConfirmation = &matches[0]
	}

	/*
		p := iCommon.NewPagination(1, 1, []iCommon.OrderBy{{Column: "rounds.number", Direction: "asc"}})
		currentRound := iService.MATCH_FILTER_STATUS_CURRENT_ROUND
		f := &iService.MatchFilter{LeagueUUID: &leagueUUID, UserUUID: &user.UUID, Status: &currentRound}
		matches, _, err := h.matchService.Get(c, p, f)
		if err != nil {
			h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
			return
		}

		var currentMatch *entities.Match
		if len(matches) > 0 {
			currentMatch = &matches[0]
		}
	*/

	groupStandings := make(map[string][]entities.GroupRoundRobinTeamResult)
	for _, f := range league.Fixtures {
		if f.Format == enums.FIXTURE_FORMAT_ROUND_ROBIN {
			sort.Slice(f.Groups, func(i, j int) bool {
				return f.Groups[i].Name[0] < f.Groups[j].Name[0]
			})
			for _, g := range f.Groups {
				standings, err := h.groupService.GetGroupRoundRobinResults(c, g.UUID)
				if err != nil {
					h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
				}
				groupStandings[g.Name] = standings
			}
		}
	}

	for _, s := range groupStandings {
		for _, r := range s {
			fmt.Printf("team: %s, matches won: %d, matches lost %d, games won %d, games lost %d \n", r.Team.Name, r.MatchWon, r.MatchLost, r.GamesWon, r.GamesLost)
		}
	}

	vmodel := &VModelPageLeagueDetail{
		User:                     user,
		League:                   league,
		Team:                     team,
		IsRatingRequired:         isRatingRequired,
		PlayersWithRatingTooHigh: playersWithRatingTooHigh,
		PlayersWithRatingTooLow:  playersWithRatingTooLow,
		// CurrentMatch:             currentMatch,
		MatchNeedsCompleteConfirmation: matchNeedsCompleteConfirmation,
		GroupStandings:                 groupStandings,
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
