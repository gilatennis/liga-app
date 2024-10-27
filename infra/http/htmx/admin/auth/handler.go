package auth

import (
	"fmt"
	"liga-app/config"
	"liga-app/enums"
	"liga-app/infra/http/htmx/common"
	iServices "liga-app/interfaces/services"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type Handler struct {
	common.BaseHandler
	cfg         *config.Config
	authService iServices.IAuthService
	userService iServices.IUserService
}

func NewHandler(
	cfg *config.Config,
	authService iServices.IAuthService,
	userService iServices.IUserService,
) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) GetLoginForm(c *gin.Context) {
	model := VModelFormLogin{}
	page := PageLogin(model)

	h.RenderView(c, page)
}

func (h *Handler) LoginWithOAuth(c *gin.Context) {
	var loginMethod enums.LOGIN_METHOD

	provider := c.Param("provider")
	values := c.Request.URL.Query()
	values.Add("provider", provider)

	redirect := c.Query("redirect")

	session := sessions.Default(c)
	session.Set("redirect", redirect)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if provider == strings.ToLower(string(enums.LOGIN_METHOD_GOOGLE)) {
		loginMethod = enums.LOGIN_METHOD_GOOGLE
	} else if provider == strings.ToLower(string(enums.LOGIN_METHOD_FACEBOOK)) {
		loginMethod = enums.LOGIN_METHOD_FACEBOOK
		values.Add("config_id", h.cfg.OAuth.FacebookConfigID)
	} else {
		// TODO: return htmx
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	c.Request.URL.RawQuery = values.Encode()

	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		user, err := h.authService.LoginWithOAuth(c, gothUser.Email, loginMethod)
		if err != nil {
			fmt.Println("SocialLogin: ", err.Error())
			// TODO: return htmx
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		isAdmin, err := h.authService.IsAdmin(c, gothUser.Email)
		if err != nil {
			fmt.Println("SocialLogin: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !isAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not admin"})
			return
		}

		session := sessions.Default(c)
		session.Set("user_uuid", user.UUID.String())
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		redirectURL := "/admin/dashboard"
		if redirect != "" {
			redirectURL = redirect
		}

		c.Redirect(http.StatusFound, redirectURL)
	} else {
		url, err := gothic.GetAuthURL(c.Writer, c.Request)
		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(c.Writer, err)
			return
		}
		// INFO: Facebook requires config_id
		// maybe create a pull request to goth to add config_id to the provider
		if loginMethod == enums.LOGIN_METHOD_FACEBOOK {
			url += "&config_id=" + h.cfg.OAuth.FacebookConfigID
		}
		http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
	}
}
