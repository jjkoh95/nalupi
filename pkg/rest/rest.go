package rest

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// New generates a new http.Server instance
func New() *http.Server {
	r := mux.NewRouter()

	// health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok!"))
	})

	r.HandleFunc("/pi/current", func(w http.ResponseWriter, r *http.Request) {

	})

	r.HandleFunc("/pi/trigger", func(w http.ResponseWriter, r *http.Request) {

	})

	// setting port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // default port
	}

	// server instance
	return &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}
}
