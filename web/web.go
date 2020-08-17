package web

import (
	"fmt"
	"github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	"grand-exchange-history/charts"
	"grand-exchange-history/item"
	"net/http"
	"sort"
	"strconv"
	"time"
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

		data, id := charts.LoadItem(id)
		fmt.Println("name:", id, data)
		var x []string
		var y []float64

		/*for k,v := range data {
			x = append(x, k)
			y = append(y, v.(float64))
		}*/

		keys := make([]string, 0, len(data))
		for k := range data {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys[len(keys)-10:] {
			sec, _ := strconv.ParseInt(k, 10, 64)
			t := time.Unix(sec/1000, 0)
			fmt.Println("sec: ", sec, t)
			x = append(x, strconv.Itoa(t.Day())+"\n"+t.Month().String())
			y = append(y, data[k].(float64))
		}

		fmt.Println(x, y)
		items := []item.ItemNameId{}
		ctx.HTML(http.StatusOK, "../../web/views/graph.html", gin.H{"title": id, "items": items, "x": x, "y": y})
	})

	router.Run(":9090")
}
