package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"reflect"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserToken struct {
	Token string `json:"token"`
}
type InvalidResponse struct {
	Error string `json:"detail"`
}

// This function will search element inside array with any type.
// Will return boolean and index for matched element.
// True and index more than 0 if element is exist.
// needle is element to search, haystack is slice of value to be search.
func InArray(needle interface{}, haystack interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

// Salt and hash
func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func (u *User) getToken() string {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	//Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["name"] = u.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	return tokenString
}

/* Set up a global string for our secret */
var mySigningKey = []byte("secret")

// JWT creation endpoint
var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	//Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["admin"] = true
	claims["name"] = "Ado Kukic"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	/* Finally, write the token to the browser window */
	w.Write([]byte(tokenString))
})

var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(5); err != nil {
		errString := fmt.Sprintf("%s: %s", http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}

	necessaryStrings := []string{"username", "password"}
	allFound := true
	var missingParam string
	formKeys := make([]string, 0, len(r.PostForm))
	for k := range r.PostForm {
		formKeys = append(formKeys, k)
	}
	for _, value := range necessaryStrings {

		stringExist, _ := InArray(value, formKeys)
		if !stringExist {
			missingParam = value
			allFound = false
		}
	}
	if !allFound {
		errString := fmt.Sprintf("%s: Missing %s param", http.StatusText(http.StatusBadRequest), missingParam)
		http.Error(w, errString, http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' ALLOW FILTERING", username)
	iter := CassandraSession.Query(query).Consistency(gocql.One).Iter()
	size := iter.NumRows()

	if size > 1 {
		errString := fmt.Sprintf("%s: Too many users found %s", http.StatusText(http.StatusInternalServerError), iter.NumRows())
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	if size == 0 {
		invalid := InvalidResponse{
			"User not found",
		}

		b, _ := json.Marshal(invalid)
		w.Write(b)
		return
	}

	m := map[string]interface{}{}
	iter.MapScan(m)

	user := User{
		m["username"].(string),
		m["password"].(string),
	}

	//salt the provided password and compare here
	if hashAndSalt([]byte(password)) == hashAndSalt([]byte(user.Password)) {
		token := UserToken{
			user.getToken(),
		}
		b, _ := json.Marshal(token)
		w.Write(b)
		return
	}

	invalid := InvalidResponse{
		"Invalid password",
	}

	b, _ := json.Marshal(invalid)
	w.Write(b)
	return

})

var IndexHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Index page should go here"))

})
