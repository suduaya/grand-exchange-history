package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path"
	"time"
	"strings"
	"fmt"
	"strconv"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/gorilla/mux"
)

const (
	host   = "127.0.0.1:8888"
	maxNum = 50
)

type router struct {
	name string
	charts.RouterOpts
}

var (
	nameItems = []string{"衬衫", "牛仔裤", "运动裤", "袜子", "冲锋衣", "羊毛衫"}
	foodItems = []string{"面包", "牛奶", "奶茶", "棒棒糖", "加多宝", "可口可乐"}

	routers = []router{
		{"bar", charts.RouterOpts{URL: host + "/bar", Text: "Bar-(柱状图)"}},
		{"kline", charts.RouterOpts{URL: host + "/kline", Text: "Kline-K 线图"}},
	}
)

func orderRouters(chartType string) []charts.RouterOpts {
	for i := 0; i < len(routers); i++ {
		if routers[i].name == chartType {
			routers[i], routers[0] = routers[0], routers[i]
			break
		}
	}

	rs := make([]charts.RouterOpts, 0)
	for i := 0; i < len(routers); i++ {
		rs = append(rs, routers[i].RouterOpts)
	}
	return rs
}

func getRenderPath(f string) string {
	return path.Join("html", f)
}

func logTracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Tracing request for %s\n", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

var seed = rand.NewSource(time.Now().UnixNano())

func randInt() []int {
	cnt := len(nameItems)
	r := make([]int, 0)
	for i := 0; i < cnt; i++ {
		r = append(r, int(seed.Int63())%maxNum)
	}
	return r
}

func genKvData() map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(nameItems); i++ {
		m[nameItems[i]] = rand.Intn(maxNum)
	}
	return m
}

func main() {
	// Avoid "404 page not found".
	router := mux.NewRouter()
	router.HandleFunc("/line/{id}", logTracing(lineHandler)).Methods("GET")
	router.HandleFunc("/summary/{id}", summary).Methods("GET")
	router.HandleFunc("/summary/contains/{id}", summaryContains).Methods("GET")

	log.Println("Run server at " + host)

	srv := &http.Server{
		Handler: router,
		Addr:    host,
		// Timeouts
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func summary(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]


	var allret []item
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
		ret := item{
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
		s, _ := strconv.ParseFloat(id, 32);

		if ret.Id == s {
			allret = append(allret, ret)

		}

	}
	fmt.Println(allret)
	json.NewEncoder(w).Encode(allret)
}

func summaryContains(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]


	var allret []item
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
		ret := item{
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
type item struct {
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
