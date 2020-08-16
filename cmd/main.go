package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"grand-exchange-history/api"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func main() {
	fmt.Println("Main is starting ...")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Postgres!")

	geApi := api.New(db)
	geApi.Start()

	//items := &item.Items{}
	//items.LoadItemsNameIds()
}
