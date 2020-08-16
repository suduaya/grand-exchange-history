package item

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	Id         float64
	LastUpdate string
	Name       string
	Buying     float64
	Selling    float64
	Margin     float64
	Overall    float64
	BuyingQtd  float64
	SellingQtd float64
	OverallQtd float64
}

func Summary(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

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
		s, _ := strconv.ParseFloat(id, 32)

		if ret.Id == s {
			allret = append(allret, ret)

		}

	}
	fmt.Println(allret)
	json.NewEncoder(w).Encode(allret)
}

func SummaryContains(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

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

		if strings.Contains(strings.ToLower(ret.Name), strings.ToLower(id)) && !v.(map[string]interface{})["members"].(bool) {
			allret = append(allret, ret)

		}

	}
	fmt.Println(allret)
	json.NewEncoder(w).Encode(allret)
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

func SelectAllItems() (ret []ItemNameId) {
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
	sqlStatement := `SELECT id, name FROM ge.items;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var name string
		var id int
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}

		item := ItemNameId{
			Id:   id,
			Name: name,
		}

		ret = append(ret, item)
	}
	return ret
}

type ItemNameId struct {
	Id   int
	Name string
}
