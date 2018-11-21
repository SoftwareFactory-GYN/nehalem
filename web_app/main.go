package main

import (
	"log"
	"net/http"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")

		if auth != "" {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/login", 301)

	})
}
func okView(w http.ResponseWriter, r *http.Request) {
	log.Println("in view")
	w.Write([]byte("OK"))
}

func main() {
	finalHandler := http.HandlerFunc(okView)

	http.Handle("/", middleware(finalHandler))
	http.ListenAndServe(":3001", nil)
}
