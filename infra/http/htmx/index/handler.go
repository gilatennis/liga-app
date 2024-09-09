package index

import (
	"liga-app/infra/http/htmx/common"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	common.BaseHandler
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Index(c *gin.Context) {
	vmodel := &VModelPageIndex{}
	page := PageIndex(vmodel)

	h.RenderView(c, page)
}
