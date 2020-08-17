package charts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/gorilla/mux"
)

func LineHandler(w http.ResponseWriter, r *http.Request) {
	page := charts.NewPage()
	params := mux.Vars(r)
	id := params["id"]
	data, id := LoadItem(id)
	page.Add(
		lineDemo(data, id),
	)
	f, err := os.Create(("line.html"))
	if err != nil {
		log.Println(err)
	}
	page.Render(w, f)
}

func lineDemo(daily map[string]interface{}, id string) *charts.Line {
	fmt.Println("Writting line")
	line := charts.NewLine()
	var ks []string
	var data []float64
	var i int = 0

	keys := make([]string, 0, len(daily))
	for k := range daily {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, v := range keys {
		data = append(data, daily[v].(float64))

		sec, _ := strconv.ParseInt(v, 10, 64)
		t := time.Unix(sec/1000, 0)
		fmt.Println("sec: ", sec, t)
		ks = append(ks, strconv.Itoa(t.Day())+"\n"+t.Month().String())

		i++

	}

	line.AddXAxis(ks[len(ks)-24:len(ks)]).AddYAxis("item_id:"+id, data[len(data)-24:len(data)],
		charts.LabelTextOpts{Show: true},
		charts.AreaStyleOpts{Opacity: 0.2},
	)
	line.SetSeriesOptions(
		charts.MLNameTypeItem{Name: "Price Average", Type: "average"},
		charts.LineOpts{Smooth: true},
		charts.MLStyleOpts{Label: charts.LabelTextOpts{Show: true, Formatter: "{a}: {b}"}},
	)
	line.SetGlobalOptions(
		charts.TitleOpts{Title: "item_id:" + id},
		charts.YAxisOpts{Name: "Price (gp)", SplitLine: charts.SplitLineOpts{Show: true}},
		charts.XAxisOpts{Name: "Date"})

	return line
}

func LoadItem(id string) (map[string]interface{}, string) {
	item_example := "http://services.runescape.com/m=itemdb_oldschool/api/graph/" + id + ".json"

	resp, err := http.Get(item_example)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var event2 map[string]interface{}
	json.Unmarshal(body, &event2)

	daily := event2["daily"].(map[string]interface{})

	return daily, id
}
