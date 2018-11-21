package main

import (
	"fmt"
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")

		if auth != "" {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/login", 301)

	})
}

func initRouter(r *mux.Router) {

	// Endpoint used to serve login page
	r.Handle("/login", LoginHandler).Methods("GET")

	//Index endpoint
	r.Handle("/", authMiddleware(middleware.JwtMiddleware.Handler(IndexHandler))).Methods("GET")

	//Static
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

}

func main() {

	port := 3001

	r := mux.NewRouter()

	// Add our routes
	initRouter(r)

	fmt.Println("Serving on: http://localhost:" + strconv.Itoa(port))

	http.ListenAndServe(":"+strconv.Itoa(port), handlers.LoggingHandler(os.Stdout, r))

}
