package duitku

import (
	"fmt"

	iService "liga-app/interfaces/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

type Handler struct {
	leagueService iService.ILeagueService
}

func NewHandler(leagueService iService.ILeagueService) *Handler {
	return &Handler{
		leagueService: leagueService,
	}
}

func (h *Handler) Callback(c *gin.Context) {
	var payload CallbackPayload
	if err := c.ShouldBind(&payload); err != nil {
		// TODO: Log error
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pp.Println("payload: ", payload)

	paymentID := fmt.Sprintf("reference=%s;publisherOrderId=%s", payload.Reference, payload.PublisherOrderID)
	paymentAmount := float64(payload.Amount)

	// TODO: do something with payment MerchantOrderID, paymentID, and paymentAmount
	_ = payload.MerchantOrderID
	_ = paymentID
	_ = paymentAmount

	txType := payload.AdditionalParam
	if txType == "league_registration" {
		leagueRegUUIDStr := payload.MerchantOrderID
		leagueRegUUID, err := uuid.Parse(leagueRegUUIDStr)
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		_, err = h.leagueService.CompleteLeagueRegistrationPayment(
			c,
			leagueRegUUID,
			paymentID,
			paymentAmount,
		)
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(200, gin.H{"message": "success"})
}
