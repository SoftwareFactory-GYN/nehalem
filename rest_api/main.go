package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var CassandraSession *gocql.Session

func initRouter(r *mux.Router) {

	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Endpoint used to serve login page
	r.Handle("/login", LoginHandler).Methods("GET")

	// Endpoint used to create new JWT's
	r.Handle("/get-token", GetTokenHandler).Methods("GET")

	//Index endpoint
	r.Handle("/", jwtMiddleware.Handler(IndexHandler)).Methods("GET")
}

func main() {

	// Cassandra init
	cassandraInit()
	CassandraSession = Session
	defer CassandraSession.Close()

	r := mux.NewRouter()

	// Add our routes
	initRouter(r)

	// Our application will run on port 3000. Here we declare the port and pass in our router.
	fmt.Println("Serving on: http://localhost:3000")
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))

}
