package web

import (
	"github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	"grand-exchange-history/item"
	"net/http"
)

func Start() {
	router := gin.Default()

	//new template engine
	router.HTMLRender = gintemplate.Default()

	router.GET("/", func(ctx *gin.Context) {
		//render only file, must full name with extension
		items := item.SelectAllItems()
		ctx.HTML(http.StatusOK, "../../web/views/page.html", gin.H{"title": "Page file title!!", "items": items})
	})

	router.Run(":9090")
}
