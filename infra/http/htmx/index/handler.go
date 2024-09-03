package index

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Index(c *gin.Context) {
	isLoggedIn := true

	session := sessions.Default(c)
	userUUIDAny := session.Get("user_uuid")
	userUUIDStr, ok := session.Get("user_uuid").(string)

	if userUUIDAny == nil || !ok {
		isLoggedIn = false
	}

	_ = userUUIDStr

	if !isLoggedIn {
		c.Redirect(http.StatusFound, "/auth/login")
		return
	} else {
		c.Redirect(http.StatusFound, "/dashboard")
	}

	c.Status(200)
}
