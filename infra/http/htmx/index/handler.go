package index

import (
	"fmt"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/infra/http/htmx/common"
	iCommon "liga-app/interfaces/common"
	"liga-app/interfaces/repositories"
	iService "liga-app/interfaces/services"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
)

type Handler struct {
	common.BaseHandler
	ruleBookService iService.IRuleBookService
	leagueService   iService.ILeagueService
	groupService    iService.IGroupService
}

func NewHandler(
	ruleBookService iService.IRuleBookService,
	leagueService iService.ILeagueService,
	groupService iService.IGroupService,
) *Handler {
	return &Handler{
		ruleBookService: ruleBookService,
		leagueService:   leagueService,
		groupService:    groupService,
	}
}

func (h *Handler) Index(c *gin.Context) {
	ruleBook, err := h.ruleBookService.GetByName(c, "main-rulebook")
	if err != nil {
		ruleBook = &entities.RuleBook{Link: ""}
	}

	user, _ := h.GetAuthenticatedUser(c)

	// INFO: for getting into a user's account, for debugging
	/*
		session := sessions.Default(c)
		session.Set("user_uuid", "37e5d5be-9d96-416b-bfa6-d6ff0a41f181")
		// pp.Println(session)

		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	*/

	pLeague := iCommon.NewPagination(1, 10, []iCommon.OrderBy{})
	season := 1
	fLeague := &repositories.LeagueFilter{Season: &season}
	leagues, _, err := h.leagueService.Get(c, pLeague, fLeague)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	groupStandings := make(map[string][]entities.GroupRoundRobinTeamResult)
	if len(leagues) > 0 {
		league := leagues[0]
		pp.Println("league.Fixtures", league.Fixtures)

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
	}
	vmodel := &VModelPageIndex{
		RuleBook:       ruleBook,
		User:           user,
		Leagues:        leagues,
		GroupStandings: groupStandings,
	}
	page := PageIndex(vmodel)

	h.RenderView(c, page)
}

func (h *Handler) GetViewGroupStandings(c *gin.Context) {
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

	view := ViewGroupStandings(groupStandings)
	h.RenderView(c, view)
}
