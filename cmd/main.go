package main

import (
	"fmt"
	"grand-exchange-history/api"
	"grand-exchange-history/item"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Main is starting ...")

	items := &item.Items{}
	items.LoadItemsNameIds()
	api.Run()

}
