package web

import (
	"fmt"
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

	router.POST("/", func(ctx *gin.Context) {
		ctx.Request.ParseForm()
		for key, value := range ctx.Request.PostForm {
			fmt.Println(key, value)
		}
		items := []item.ItemNameId{}
		ctx.HTML(http.StatusOK, "../../web/views/page.html", gin.H{"title": "Page file title!!", "items": items})
	})

	router.GET("/graph/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		fmt.Println("name:", id)
		items := []item.ItemNameId{}
		ctx.HTML(http.StatusOK, "../../web/views/graph.html", gin.H{"title": id, "items": items})
	})

	router.Run(":9090")
}
