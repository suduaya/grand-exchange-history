package item

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Item struct {
	Id         float64
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
