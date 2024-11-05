package middlewares

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// FIXME: the logic is incorrect
func NewHistoryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if c.Request.Method == "GET" {
			history := []string{}
			existingHistory := session.Get("history")
			if existingHistory != nil {
				history = existingHistory.([]string)
			}

			history = append(history, c.Request.URL.Path)

			session.Set("history", history)

			err := session.Save()
			if err != nil {
				fmt.Print("Error: history middleware: session.Save(): ", err)
				return
			}
		}
	}
}
