package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"grand-exchange-history/item"
	"net/http"
)

func (t *API) ItemSearchHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	item_name := params["item_name"]
	json.NewEncoder(w).Encode(item.ItemNameContains(item_name))
}
