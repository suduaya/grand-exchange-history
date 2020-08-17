package api

import (
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type API struct {
	mux    *mux.Router
	server *http.Server
	db     *sql.DB
	//itemMgr    types.ItemPrv
}

// Middleware
func logTracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Tracing request for %s\n", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

func New(db *sql.DB) *API {
	// Avoid "404 page not found".
	router := mux.NewRouter()

	c := make(chan struct{}, 100) // max requests example

	router.Use(func(next http.Handler) http.Handler {
		// Limiting the degree of concurrency.
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Counting semaphore using a buffered channel.
			select {
			case c <- struct{}{}:
				defer func() { <-c }()

				// Call the next handler, which can be another middleware in the chain, or the final handler.
				next.ServeHTTP(w, r)
			default:
				w.WriteHeader(http.StatusTooManyRequests)
			}
		})
	},
		func(next http.Handler) http.Handler {
			// Manipulate the header for all the HTTP(S) responses.
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")

				// Call the next handler, which can be another middleware in the chain, or the final handler.
				next.ServeHTTP(w, r)
			})
		})

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8888",
		// Timeouts
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	api := &API{
		mux:    router,
		server: srv,
		db:     db,
	}

	router.HandleFunc("/search/{item_name}", logTracing(api.ItemSearchHandler)).Methods("GET")

	log.Println("Run server at " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
	return api
}

func (t *API) Start() error {
	return t.server.ListenAndServe()
}

// Shutdown attempts to close the http server.
func (t *API) Close() error {
	return t.server.Shutdown(context.Background())
}
