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
		// 745a8a89-8516-4096-8c02-948c4c681d45 // girigyat@gmail.com
		session.Set("user_uuid", "745a8a89-8516-4096-8c02-948c4c681d45")

		// 57de439c-1acb-489e-ae0c-6dc3a1b8ee67 // pop the cherry captain
		// session.Set("user_uuid", "57de439c-1acb-489e-ae0c-6dc3a1b8ee67")

		// jody
		// session.Set("user_uuid", "d9529ebd-791b-4f54-81da-ee684f8c26d6")

		// fadel
		// session.Set("user_uuid", "7228c63c-70f4-48c9-ac1b-0ee07d18b47a")
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
