package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world!")
}

func main() {

	port := 8000

	http.Handle("/", middlewareOne(http.HandlerFunc(handler)))

	log.Println(fmt.Sprintf("Server starting: http://localhost:%d", port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
