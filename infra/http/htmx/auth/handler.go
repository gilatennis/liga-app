package auth

import (
	"errors"
	"fmt"
	"liga-app/config"
	"liga-app/enums"
	"liga-app/infra/http/htmx/common/renderer"
	iServices "liga-app/interfaces/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gorillaSessions "github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	gothGoogle "github.com/markbates/goth/providers/google"
)

type Handler struct {
	cfg         *config.Config
	authService iServices.IAuthService
}

func NewHandler(cfg *config.Config, authService iServices.IAuthService) *Handler {
	SetupSocialLoginProviders(cfg)
	return &Handler{
		cfg:         cfg,
		authService: authService,
	}
}

func SetupSocialLoginProviders(cfg *config.Config) {
	store := gorillaSessions.NewCookieStore([]byte("secret"))
	gothic.Store = store
	goth.UseProviders(
		gothGoogle.New(cfg.OAuth.GoogleClientID, cfg.OAuth.GoogleSecret, cfg.OAuth.GoogleCallbackURL),
		facebook.New(cfg.OAuth.FacebookAppID, cfg.OAuth.FacebookSecret, cfg.OAuth.FacebookCallbackURL),
	)
}

func (h *Handler) GetRegisterForm(c *gin.Context) {
	redirect := c.Query("redirect")
	fmt.Println("DEBUG: [GetRegisterForm] redirect: ", redirect)

	model := VModelFormRegister{}
	page := PageRegister(model, redirect)

	r := renderer.New(c.Request.Context(), http.StatusOK, page)
	c.Render(http.StatusOK, r)
}

func (h *Handler) SubmitRegisterForm(c *gin.Context) {
	redirect := c.Query("redirect")
	var model VModelFormRegister
	if err := c.ShouldBind(&model); err != nil {
		// err = fmt.Errorf("%s: %w", logTag, err)

		log.Println(err)
		model.ServerError = err.Error()

		content := ViewFormRegister(model, redirect)
		r := renderer.New(c.Request.Context(), http.StatusInternalServerError, content)
		c.Render(http.StatusInternalServerError, r)
		return
	}

	if !model.Validate() {
		content := ViewFormRegister(model, redirect)
		r := renderer.New(c.Request.Context(), http.StatusBadRequest, content)
		c.Render(http.StatusBadRequest, r)
		return
	}

	var magicLinkRedirect *string = nil
	if redirect != "" {
		magicLinkRedirect = &redirect
	}
	err := h.authService.RegisterWithMagicLink(c, model.Email, magicLinkRedirect)
	if err != nil {
		model.ServerError = errors.Unwrap(err).Error()
		content := ViewFormRegister(model, redirect)
		r := renderer.New(c.Request.Context(), http.StatusInternalServerError, content)
		c.Render(http.StatusInternalServerError, r)
		return
	}

	msg := fmt.Sprintf("Please check your email %s to verify your account", model.Email)
	content := ViewRegisterSuccess(msg)
	r := renderer.New(c.Request.Context(), http.StatusOK, content)
	c.Render(http.StatusOK, r)
}

func (h *Handler) CompleteMagicLinkRegistration(c *gin.Context) {
	redirect := c.Query("redirect")
	email := c.Query("email")
	token := c.Query("token")

	if email == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and token are required"})
		return
	}

	user, err := h.authService.CompleteMagicLinkRegistration(c, email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	session.Set("user_uuid", user.UUID.String())
	// FIXME: setting max age makes session not working after redirect
	// session.Options(sessions.Options{MaxAge: h.cfg.Session.MaxAge})
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	redirectURL := "/dashboard"
	if redirect != "" {
		redirectURL = redirect
	}
	c.Redirect(http.StatusFound, redirectURL)
}

func (h *Handler) GetLoginForm(c *gin.Context) {
	model := VModelFormLogin{}
	page := PageLogin(model)

	r := renderer.New(c.Request.Context(), http.StatusOK, page)
	c.Render(http.StatusOK, r)
}

func (h *Handler) SubmitLoginForm(c *gin.Context) {
	redirect := c.Query("redirect")
	var model VModelFormLogin
	if err := c.ShouldBind(&model); err != nil {
		// err = fmt.Errorf("%s: %w", logTag, err)

		log.Println(err)
		model.ServerError = err.Error()

		content := LoginForm(model, redirect)
		r := renderer.New(c.Request.Context(), http.StatusInternalServerError, content)
		c.Render(http.StatusInternalServerError, r)
		return
	}

	if !model.Validate() {
		content := LoginForm(model, redirect)
		r := renderer.New(c.Request.Context(), http.StatusBadRequest, content)
		c.Render(http.StatusBadRequest, r)
		return
	}

	var magicLinkRedirect *string = nil
	if redirect != "" {
		magicLinkRedirect = &redirect
	}
	err := h.authService.LoginWithMagicLink(c, model.Email, magicLinkRedirect)
	if err != nil {
		model.ServerError = errors.Unwrap(err).Error()
		content := LoginForm(model, redirect)
		r := renderer.New(c.Request.Context(), http.StatusInternalServerError, content)
		c.Render(http.StatusInternalServerError, r)
		return
	}

	msg := fmt.Sprintf("Please check your email %s for login link", model.Email)
	content := loginSuccess(msg)
	r := renderer.New(c.Request.Context(), http.StatusOK, content)
	c.Render(http.StatusOK, r)
}

func (h *Handler) LoginWithMagicLink(c *gin.Context) {
	redirect := c.Query("redirect")
	email := c.Query("email")
	token := c.Query("token")

	if email == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and token are required"})
		return
	}

	user, err := h.authService.CompleteMagicLinkLogin(c, email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	session.Set("user_uuid", user.UUID.String())
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	redirectURL := "/dashboard"
	if redirect != "" {
		redirectURL = redirect
	}
	c.Redirect(http.StatusFound, redirectURL)
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

		session := sessions.Default(c)
		session.Set("user_uuid", user.UUID.String())
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		redirectURL := "/dashboard"
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

func (h *Handler) LoginWithOAuthCallback(c *gin.Context) {
	var loginMethod enums.LOGIN_METHOD

	provider := c.Param("provider")
	values := c.Request.URL.Query()
	values.Add("provider", provider)
	c.Request.URL.RawQuery = values.Encode()

	if provider == strings.ToLower(string(enums.LOGIN_METHOD_GOOGLE)) {
		loginMethod = enums.LOGIN_METHOD_GOOGLE
	} else if provider == strings.ToLower(string(enums.LOGIN_METHOD_FACEBOOK)) {
		loginMethod = enums.LOGIN_METHOD_FACEBOOK
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Println("SocialLoginCallback-CompleteUserAuth: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.LoginWithOAuth(c, gothUser.Email, loginMethod)
	if err != nil {
		fmt.Println("SocialLoginCallback-SocialLogin: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	session.Set("user_uuid", user.UUID.String())
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	redirectURL := "/dashboard"

	redirectInSession := sessions.Default(c).Get("redirect")
	if redirectInSession != nil && redirectInSession.(string) != "" {
		redirectURL = redirectInSession.(string)
	}
	// clear redirect in session
	session.Delete("redirect")

	fmt.Println("DEBUG: [SocialLoginCallback] redirectURL: ", redirectURL)
	c.Redirect(http.StatusFound, redirectURL)
}

func (h *Handler) LogOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/auth/login")
}
