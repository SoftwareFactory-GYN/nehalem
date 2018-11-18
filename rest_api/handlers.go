package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"reflect"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type InvalidResponse struct {
	TypeOfError string `json:"type"`
	Description string `json:"detail"`
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

	w.Write([]byte("Login page should go here"))
})

var IndexHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Index page should go here"))

})
