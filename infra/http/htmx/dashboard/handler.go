package dashboard

import (
	"liga-app/infra/http/htmx/common"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
)

type Handler struct {
	common.BaseHandler
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Index(c *gin.Context) {
	user, _ := h.GetAuthenticatedUser(c)
	pp.Println(user)
}
