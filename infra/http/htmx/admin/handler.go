package admin

import "github.com/gin-gonic/gin"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Index(c *gin.Context) {
}
