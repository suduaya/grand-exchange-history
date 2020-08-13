package api

import (
	"github.com/gorilla/mux"
	"grand-exchange-history/charts"
	"grand-exchange-history/item"
	"log"
	"net/http"
	"time"
)

// Middleware
func logTracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Tracing request for %s\n", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

func Run() {
	// Avoid "404 page not found".
	router := mux.NewRouter()

	router.HandleFunc("/line/{id}", logTracing(charts.LineHandler)).Methods("GET")
	router.HandleFunc("/summary/{id}", logTracing(item.Summary)).Methods("GET")
	router.HandleFunc("/summary/contains/{id}", logTracing(item.SummaryContains)).Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8888",
		// Timeouts
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	log.Println("Run server at " + srv.Addr)
	log.Fatal(srv.ListenAndServe())

}
