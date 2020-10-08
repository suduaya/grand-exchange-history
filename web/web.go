package web

import (
	"fmt"
	"grand-exchange-history/item"
	"net/http"
	"strconv"
	"strings"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
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
		var item_name string
		for key, value := range ctx.Request.PostForm {
			fmt.Println(key, value)
			for _, v := range value {
				item_name += v
			}
		}

		items := item.SelectAllItems()
		var ret_items []item.ItemNameId
		for _, item := range items {
			if strings.Contains(strings.ToLower(item.Name), strings.ToLower(item_name)) {
				ret_items = append(ret_items, item)
			}
		}
		ctx.HTML(http.StatusOK, "../../web/views/page.html", gin.H{"title": "Page file title!!", "items": ret_items})
	})

	router.GET("/graph/:id/weekly", func(ctx *gin.Context) {
		id := ctx.Param("id")

		s, _ := strconv.ParseFloat(id, 64)
		item_struct := item.GetItemById(s)[0]
		fmt.Println("item_struct: ", item_struct)

		x, y := item.LoadItemPriceHistory(id, 7)
		items := []item.ItemNameId{}
		ctx.HTML(http.StatusOK, "../../web/views/graph.html", gin.H{"id": id, "items": items, "x": x, "y": y, "item_struct": item_struct})
	})

	router.GET("/graph/:id/monthly", func(ctx *gin.Context) {
		id := ctx.Param("id")

		s, _ := strconv.ParseFloat(id, 64)
		item_struct := item.GetItemById(s)[0]
		fmt.Println("item_struct: ", item_struct)

		x, y := item.LoadItemPriceHistory(id, 30)
		items := []item.ItemNameId{}
		ctx.HTML(http.StatusOK, "../../web/views/graph.html", gin.H{"id": id, "items": items, "x": x, "y": y, "item_struct": item_struct})
	})

	router.GET("/graph/:id/quarter", func(ctx *gin.Context) {
		id := ctx.Param("id")

		s, _ := strconv.ParseFloat(id, 64)
		item_struct := item.GetItemById(s)[0]
		fmt.Println("item_struct: ", item_struct)

		x, y := item.LoadItemPriceHistory(id, 90)
		items := []item.ItemNameId{}
		ctx.HTML(http.StatusOK, "../../web/views/graph.html", gin.H{"id": id, "items": items, "x": x, "y": y, "item_struct": item_struct})
	})

	router.Run(":9090")
}
