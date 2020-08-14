package item

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

type Items struct {
	Items []Item
	Db    sql.DB
}

type Item struct {
	Id         float64 `json:"id"`
	LastUpdate string  `json:"last_update"`
	Name       string  `json:"name"`
	Buying     float64 `json:"buying_at"`
	Selling    float64 `json:"selling_at"`
	Margin     float64 `json:"margin"`
	Overall    float64 `json:"overall"`
	BuyingQtd  float64 `json:"buying_qtd"`
	SellingQtd float64 `json:"selling_qtd"`
	OverallQtd float64 `json:"overall_qtd"`
}

func ItemNameContains(item_name string) []Item {
	var allret []Item
	api_endpoint := "https://rsbuddy.com/exchange/summary.json"

	resp, err := http.Get(api_endpoint)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var event map[string]interface{}
	json.Unmarshal(body, &event)
	for _, v := range event {
		ret := Item{
			Id:         v.(map[string]interface{})["id"].(float64),
			Name:       v.(map[string]interface{})["name"].(string),
			Buying:     v.(map[string]interface{})["buy_average"].(float64),
			Selling:    v.(map[string]interface{})["sell_average"].(float64),
			Margin:     v.(map[string]interface{})["buy_average"].(float64) - v.(map[string]interface{})["sell_average"].(float64),
			Overall:    v.(map[string]interface{})["overall_average"].(float64),
			BuyingQtd:  v.(map[string]interface{})["buy_quantity"].(float64),
			SellingQtd: v.(map[string]interface{})["sell_quantity"].(float64),
			OverallQtd: v.(map[string]interface{})["overall_quantity"].(float64),
		}

		if strings.Contains(strings.ToLower(ret.Name), strings.ToLower(item_name)) {
			allret = append(allret, ret)
		}

	}
	fmt.Println(allret)
	return allret
}

func (items Items) LoadItemsNameIds() {
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

	api_endpoint := "https://rsbuddy.com/exchange/summary.json"

	resp, err := http.Get(api_endpoint)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var event map[string]interface{}
	json.Unmarshal(body, &event)
	for _, v := range event {
		id := v.(map[string]interface{})["id"].(float64)
		name := v.(map[string]interface{})["name"].(string)
		insertItem(db, id, name)
	}

}

func insertItem(db *sql.DB, id float64, name string) {
	fmt.Println("Inserting:", id, name)
	sqlStatement := `INSERT INTO ge.items(id, name) VALUES ($1, $2);`
	_, err := db.Exec(sqlStatement, id, name)
	if err != nil {
		fmt.Println(err)
	}
}
