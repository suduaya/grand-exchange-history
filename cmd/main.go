package main

import (
	"fmt"
	"grand-exchange-history/web"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Main is starting ...")

	//items := &item.Items{}
	//items.LoadItemsNameIds()
	//go api.Run()
	web.Start()

}
