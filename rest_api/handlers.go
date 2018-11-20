package main

import (
	"encoding/json"
	"fmt"
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/secret"
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/user"
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/utils"
	"log"
	"net/http"
)

var mySigningKey = secret.GetSigningKey()

type InvalidResponse struct {
	Error string `json:"detail"`
}

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

		stringExist, _ := utils.InArray(value, formKeys)
		if !stringExist {
			missingParam = value
			allFound = false
		}
	}
	if !allFound {
		errString := fmt.Sprintf("%s: missing %s param", http.StatusText(http.StatusBadRequest), missingParam)
		http.Error(w, errString, http.StatusBadRequest)
		return
	}

	attemptingUser := user.User{
		"",
		r.PostForm.Get("username"),
		r.PostForm.Get("password"),
	}

	if !attemptingUser.Exists() {
		invalid := InvalidResponse{
			"user not found",
		}

		b, _ := json.Marshal(invalid)
		w.Write(b)
		return
	}

	existingUser, err := user.FetchUser(attemptingUser.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//salt the provided password and compare here
	if user.ComparePasswords(existingUser.Password, []byte(attemptingUser.Password)) {
		token := user.UserToken{
			existingUser.GetToken(),
		}
		b, _ := json.Marshal(token)
		w.Write(b)
		return
	}

	invalid := InvalidResponse{
		"invalid password",
	}

	b, _ := json.Marshal(invalid)
	w.Write(b)
	return

})

var RegisterHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

		stringExist, _ := utils.InArray(value, formKeys)
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

	newUser := user.User{
		"",
		r.PostForm.Get("username"),
		r.PostForm.Get("password"),
	}

	if newUser.Exists() {
		invalid := InvalidResponse{
			"User already exists",
		}

		b, _ := json.Marshal(invalid)
		w.Write(b)
		return
	}

	newUser.Create()

	token := user.UserToken{
		newUser.GetToken(),
	}
	b, _ := json.Marshal(token)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
	return

})

var IndexHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Index page should go here"))

})
