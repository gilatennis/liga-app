package common

import (
	"fmt"
	"liga-app/entities"
	"liga-app/infra/http/htmx/common/renderer"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BaseHandler struct{}

func (h *BaseHandler) RenderView(c *gin.Context, view templ.Component) {
	r := renderer.New(c, http.StatusOK, view)
	c.Render(http.StatusOK, r)
}

func (h *BaseHandler) RenderViewWithStatus(c *gin.Context, view templ.Component, status int) {
	r := renderer.New(c, status, view)
	c.Render(http.StatusOK, r)
}

func (h *BaseHandler) GetAuthenticatedUser(c *gin.Context) (*entities.User, error) {
	userAny, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return nil, fmt.Errorf("unauthorized")
	}
	user, ok := userAny.(*entities.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return nil, fmt.Errorf("unauthorized")
	}

	return user, nil
}

func (h *BaseHandler) GetURLUUIDParams(c *gin.Context, params []string) ([]uuid.UUID, error) {
	var values []uuid.UUID
	for _, param := range params {
		value := c.Param(param)
		if value == "" {
			err := fmt.Errorf("GetURLParams: param %s not found", param)
			return nil, err
		}
		valUUID, err := uuid.Parse(value)
		if err != nil {
			err = fmt.Errorf("GetURLParams: uuid.Parse error: %w", err)
			return nil, err
		}
		values = append(values, valUUID)
	}

	return values, nil
}

func (h *BaseHandler) GetQueryUUIDParams(c *gin.Context, params []string) ([]uuid.UUID, error) {
	var values []uuid.UUID
	for _, param := range params {
		value := c.Query(param)
		if value == "" {
			err := fmt.Errorf("GetQueryUUIDParams: param %s not found", param)
			return nil, err
		}
		valUUID, err := uuid.Parse(value)
		if err != nil {
			err = fmt.Errorf("GetQueryUUIDParams: uuid.Parse error: %w", err)
			return nil, err
		}
		values = append(values, valUUID)
	}

	return values, nil
}

func (b *BaseHandler) RenderErrorPage(c *gin.Context, status int, message string) {
	page := PageError(status, message)
	r := renderer.New(c, status, page)
	c.Render(status, r)
}

func (h *BaseHandler) RenderResponseCreatedPage(c *gin.Context, location string, page templ.Component, message string) error {
	alert := NewAlertPayload(AlertTypeSuccess, message)
	jsonAlertMessage, err := alert.ToJSON()
	if err != nil {
		err = fmt.Errorf("htmx.common.base_handler.RenderResponseCreatedPage() : %w", err)
		return err
	}

	c.Header("HX-Trigger-After-Settle", string(jsonAlertMessage))
	c.Header("HX-Push-Url", location)

	r := renderer.New(c.Request.Context(), http.StatusOK, page)
	c.Render(http.StatusCreated, r)

	return nil
}

/**
* ResponseSuccessPushURL is a helper function to send a success message to the client
* location is the URL that will be pushed on client side via htmx
* message is the message that will be displayed on client side
* this will not redirect server side.
* so you need to handle the page to render as a response
* */
func (h *BaseHandler) ResponseSuccessPushURL(c *gin.Context, message string, location string) {
	alert := NewAlertPayload(AlertTypeSuccess, message)
	jsonAlertMessage, err := alert.ToJSON()
	if err != nil {
		err = fmt.Errorf("htmx.common.base_handler.ResponseSuccessPushURL() : %w", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("HX-Trigger-After-Settle", string(jsonAlertMessage))
	c.Header("HX-Push-Url", location)
}

/**
* Alert success message and redirect client side
* */
func (h *BaseHandler) ResponseSuccessRedirect(c *gin.Context, message, location string) error {
	alert := AlertPayloadWithRedirect{
		Alert: AlertPayloadWithRedirectDetail{
			Type:     AlertTypeSuccess,
			Message:  message,
			Redirect: location,
		},
	}

	jsonAlertMessage, err := alert.ToJSON()
	if err != nil {
		return err
	}

	c.Header("HX-Trigger", string(jsonAlertMessage))

	return nil
}

func (h *BaseHandler) ResponseAlertError(c *gin.Context, message string) error {
	alert := NewAlertPayload(AlertTypeError, message)

	jsonAlertMessage, err := alert.ToJSON()
	if err != nil {
		err = fmt.Errorf("htmx.common.base_handler.RenderResponseCreatedPage() : %w", err)
		return err
	}

	c.Header("HX-Trigger-After-Settle", string(jsonAlertMessage))
	return nil
}

func (h *BaseHandler) ResponseAlertSuccess(c *gin.Context, message string) error {
	alert := NewAlertPayload(AlertTypeSuccess, message)

	jsonAlertMessage, err := alert.ToJSON()
	if err != nil {
		err = fmt.Errorf("htmx.common.base_handler.RenderResponseCreatedPage() : %w", err)
		return err
	}

	c.Header("HX-Trigger-After-Settle", string(jsonAlertMessage))
	return nil
}
