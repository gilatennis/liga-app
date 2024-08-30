package duitku

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
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

	c.JSON(200, gin.H{"message": "success"})
}
