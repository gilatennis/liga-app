package index

import (
	"liga-app/entities"
	"liga-app/infra/http/htmx/common"
	iService "liga-app/interfaces/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	common.BaseHandler
	ruleBookService iService.IRuleBookService
}

func NewHandler(ruleBookService iService.IRuleBookService) *Handler {
	return &Handler{ruleBookService: ruleBookService}
}

func (h *Handler) Index(c *gin.Context) {
	ruleBook, err := h.ruleBookService.GetByName(c, "main-rulebook")
	if err != nil {
		ruleBook = &entities.RuleBook{Link: ""}
	}

	// INFO: for getting into a user's account, for debugging
	/*
		session := sessions.Default(c)
		// session.Set("user_uuid", "d9529ebd-791b-4f54-81da-ee684f8c26d6")
		session.Set("user_uuid", "7228c63c-70f4-48c9-ac1b-0ee07d18b47a")
		pp.Println(session)

		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	*/

	vmodel := &VModelPageIndex{
		RuleBook: ruleBook,
	}
	page := PageIndex(vmodel)

	h.RenderView(c, page)
}
