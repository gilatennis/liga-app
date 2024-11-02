package player

import (
	"liga-app/entities"
	"liga-app/infra/http/htmx/common"
	"liga-app/infra/http/htmx/common/components"
	"liga-app/interfaces/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
)

type Handler struct {
	common.BaseHandler
	userService     services.IUserService
	userInfoService services.IUserInfoService
}

func NewHandler(
	userService services.IUserService,
	userInfoService services.IUserInfoService,
) *Handler {
	return &Handler{
		userService:     userService,
		userInfoService: userInfoService,
	}
}

func (h *Handler) GetPlayerDetailPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"player_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	playerUUID := params[0]

	player, err := h.userService.GetUserByUUID(c, playerUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}
	pp.Println(player.Ratings)

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

	vmodel := VModelPagePlayerDetail{
		User:   user,
		Player: player,
		VModelFormPlayer: &VModelFormPlayer{
			ProfilePictureCurrent: player.ProfilePictureURL,
			ProvinceOptions:       provinceOptions,
			Name:                  player.Name,
			Phone:                 player.Phone,
			Province:              player.Province,
			RatingSeason1:         player.GetRatingBySeason(1),
		},
	}

	page := PagePlayerDetail(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) PutPlayerDetail(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"player_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}

	playerUUID := params[0]

	player, err := h.userService.GetUserByUUID(c, playerUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
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

	pp.Println(player)

	vmodel := VModelPagePlayerDetail{
		User:   user,
		Player: player,
		VModelFormPlayer: &VModelFormPlayer{
			ProfilePictureCurrent: player.ProfilePictureURL,
			ProvinceOptions:       provinceOptions,
			Name:                  player.Name,
			Phone:                 player.Phone,
			Province:              player.Province,
			RatingSeason1:         player.GetRatingBySeason(1),
		},
	}

	err = c.Bind(vmodel.VModelFormPlayer)
	if err != nil {
		view := ViewFormEditProfile(*vmodel.VModelFormPlayer, player.UUID)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	if !vmodel.VModelFormPlayer.Validate() {
		view := ViewFormEditProfile(*vmodel.VModelFormPlayer, player.UUID)
		h.RenderViewWithStatus(c, view, http.StatusBadRequest)
		return
	}

	player.Name = vmodel.VModelFormPlayer.Name
	player.Phone = vmodel.VModelFormPlayer.Phone
	player.Province = vmodel.VModelFormPlayer.Province

	if vmodel.VModelFormPlayer.RatingSeason1 != nil {
		if len(player.Ratings) == 0 {
			player.Ratings = append(player.Ratings, entities.UserRating{
				UserUUID:   player.UUID,
				NTRPRating: *vmodel.VModelFormPlayer.RatingSeason1,
				Season:     1,
			})
		} else {
			for i, r := range player.Ratings {
				if r.Season == 1 {
					r.NTRPRating = *vmodel.VModelFormPlayer.RatingSeason1
				}
				player.Ratings[i] = r
			}
		}
	}

	pp.Println(player)
	_, err = h.userService.UpdateUser(c, user.UUID, player)
	if err != nil {
		vmodel.VModelFormPlayer.base.ServerError = err.Error()
		view := ViewFormEditProfile(*vmodel.VModelFormPlayer, player.UUID)
		h.RenderViewWithStatus(c, view, http.StatusInternalServerError)
		return
	}

	view := ViewFormEditProfile(*vmodel.VModelFormPlayer, player.UUID)
	h.ResponseAlertSuccess(c, "Informasi user berhasil di update")
	h.RenderView(c, view)
}
