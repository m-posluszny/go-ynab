package panel

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-posluszny/go-ynab/src/auth"
	"github.com/m-posluszny/go-ynab/src/db"
)

func SetRoutes(r *gin.Engine) {
	userPanel := r.Group("/panel")
	userPanel.Use(auth.AuthRequired)
	{
		userPanel.GET("/", func(c *gin.Context) {
			uid := auth.GetUIDFromSession(c)
			RenderPanel(c, uid, http.StatusOK)
		})

	}
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/panel")
	})

}

func RenderPanel(c *gin.Context, uid string, status int) {
	dbx := db.GetDbRead()
	creds, err := auth.GetUserFromUid(dbx, uid)
	if err != nil {
		panic(err)
	}
	c.HTML(status, "panel.html", gin.H{"username": creds.Username})
}
