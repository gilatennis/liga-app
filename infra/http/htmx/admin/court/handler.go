package court

import (
	"liga-app/infra/http/htmx/common"
	iCommon "liga-app/interfaces/common"
	"liga-app/interfaces/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	common.BaseHandler
	courtService services.ICourtService
}

func NewHandler(courtService services.ICourtService) *Handler {
	return &Handler{
		courtService: courtService,
	}
}

func (h *Handler) GetCourtListPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	p := iCommon.NewPagination(1, 100, []iCommon.OrderBy{})
	courts, err := h.courtService.Get(c, p)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageCourtList{
		User:   user,
		Courts: courts,
	}

	page := PageCourtList(vmodel)
	h.RenderView(c, page)
}

func (h *Handler) GetCourtDetailPage(c *gin.Context) {
	user, err := h.GetAuthenticatedUser(c)
	if err != nil {
		h.RenderErrorPage(c, http.StatusUnauthorized, err.Error())
		return
	}

	params, err := h.GetURLUUIDParams(c, []string{"court_uuid"})
	if err != nil {
		h.RenderErrorPage(c, http.StatusBadRequest, err.Error())
		return
	}
	courtUUID := params[0]

	court, err := h.courtService.GetByUUID(c, courtUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	availableDays, err := h.courtService.GetAdminCourtAvailableDays(c, courtUUID)
	if err != nil {
		h.RenderErrorPage(c, http.StatusInternalServerError, err.Error())
		return
	}

	vmodel := VModelPageCourtDetail{
		User:          user,
		Court:         court,
		AvailableDays: availableDays,
	}

	page := PageCourtDetail(vmodel)
	h.RenderView(c, page)
}
