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
		session.Set("user_uuid", "0fe90a49-f1db-4ece-80d7-4227867bc3f1")
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
