package profile

import (
	"liga-app/infra/http/htmx/common"
	"liga-app/infra/http/htmx/common/components"
	"mime/multipart"
	"net/http"

	iService "liga-app/interfaces/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	common.BaseHandler
	userInfoService iService.IUserInfoService
}

func NewHandler(userInfoService iService.IUserInfoService) *Handler {
	return &Handler{
		userInfoService: userInfoService,
	}
}

func (h *Handler) GetProfileDetailPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	vmodel := VModelPageProfileDetail{
		User: user,
	}
	page := PageProfileDetail(&vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetProfileEditPage(c *gin.Context) {
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

	vmodel := VModelPageProfileEdit{
		User: user,
		VModelFormProfile: &VModelFormProfile{
			ProfilePictureCurrent: user.ProfilePictureURL,
			ProvinceOptions:       provinceOptions,
			Name:                  user.Name,
			Phone:                 user.Phone,
			Province:              user.Province,
		},
	}

	page := PageProfileEdit(&vmodel)
	h.RenderView(c, page)
}

func (h *Handler) PutUpdateProfile(c *gin.Context) {
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

	form := VModelFormProfile{
		ProfilePictureCurrent: user.ProfilePictureURL,
		ProvinceOptions:       provinceOptions,
		Gender:                string(user.Gender),
	}

	err = c.Bind(&form)
	if err != nil {
		form.base.ServerError = err.Error()
		view := ViewFormEditProfile(&form)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	if !form.Validate() {
		view := ViewFormEditProfile(&form)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	fileHeader, err := c.FormFile("profile_picture")
	if err != nil && err.Error() != "http: no such file" {
		form.base.ServerError = err.Error()
		view := ViewFormEditProfile(&form)
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	var file *multipart.File = nil

	if fileHeader != nil {
		_file, err := fileHeader.Open()
		if err != nil {
			form.base.ServerError = err.Error()
			view := ViewFormEditProfile(&form)
			h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
			return
		}
		file = &_file
	}

	user.Name = form.Name
	user.Phone = form.Phone
	user.Province = form.Province

	_, err = h.userInfoService.SaveUserInfo(c, user, file)
	if err != nil {
		form.base.ServerError = err.Error()
		view := ViewFormEditProfile(&form)
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	h.ResponseSuccessRedirect(c, "informasi user berhasil di update", "/dashboard/profile")
}
