package dashboard

import (
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/infra/http/htmx/common"
	"liga-app/infra/http/htmx/common/components"
	"net/http"

	iCommon "liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"
	iService "liga-app/interfaces/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	common.BaseHandler
	userInfoService iService.IUserInfoService
	leagueService   iService.ILeagueService
}

func NewHandler(userInfoService iService.IUserInfoService, leagueService iService.ILeagueService) *Handler {
	return &Handler{
		userInfoService: userInfoService,
		leagueService:   leagueService,
	}
}

func (h *Handler) Index(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	if !h.isUserInfoComplete(user) {
		c.Redirect(http.StatusTemporaryRedirect, "/dashboard/complete-user-info")
		return
	}

	p := iCommon.NewPagination(1, 6, []iCommon.OrderBy{{Column: "created_at", Direction: "desc"}})
	f := &iRepo.LeagueFilter{UserUUID: &user.UUID}
	leagues, _, err := h.leagueService.Get(c, p, f)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := &VModelPageDashboardIndex{
		User:    user,
		Leagues: leagues,
	}
	page := PageDashboardIndex(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) isUserInfoComplete(user *entities.User) bool {
	return user.Name != "" &&
		user.Phone != "" &&
		user.Gender != "" &&
		user.Province != ""
}

func (h *Handler) GetCompleteUserInfoPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	provinces, err := h.userInfoService.GetProvinceOptions(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	provinceOptions := []components.InputSelectOption{{Value: "", Text: "-"}}
	for _, c := range provinces {
		provinceOptions = append(provinceOptions, components.InputSelectOption{
			Value: c.IDProvince,
			Text:  c.Name,
		})
	}

	vmodel := &VModelPageUserInfo{
		User: user,
		VModelFormUserInfo: &VModelFormUserInfo{
			ProvinceOptions: provinceOptions,
			Province:        "31", // set default value to Jakarta
		},
	}
	page := PageCompleteUserInfo(vmodel)

	h.RenderView(c, page)
}

func (h *Handler) PostSubmitUserInfoForm(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	provinces, err := h.userInfoService.GetProvinceOptions(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	provinceOptions := []components.InputSelectOption{{Value: "", Text: "-"}}
	for _, c := range provinces {
		provinceOptions = append(provinceOptions, components.InputSelectOption{
			Value: c.IDProvince,
			Text:  c.Name,
		})
	}

	vmodelForm := VModelFormUserInfo{
		ProvinceOptions: provinceOptions,
	}
	err = c.Bind(&vmodelForm)
	if err != nil {
		view := ViewFormCompleteUserInfo(&vmodelForm)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	if !vmodelForm.Validate() {
		view := ViewFormCompleteUserInfo(&vmodelForm)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	user.Name = vmodelForm.Name
	user.Phone = vmodelForm.Phone
	user.Province = vmodelForm.Province
	user.Gender = enums.GENDER(vmodelForm.Gender)

	_, err = h.userInfoService.SaveUserInfo(c, user)
	if err != nil {
		vmodelForm.base.ServerError = err.Error()
		view := ViewFormCompleteUserInfo(&vmodelForm)
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	view := ViewCompleteUserInfoSuccess()
	h.RenderView(c, view)
}
