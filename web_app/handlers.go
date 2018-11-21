package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
)

type InvalidResponse struct {
	Error string `json:"detail"`
}

var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	fp := path.Join("templates", "login.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Println(err)
		res := InvalidResponse{
			err.Error(),
		}

		b, _ := json.Marshal(res)
		w.Write(b)
		return

	}

	if err := tmpl.Execute(w, nil); err != nil {

		log.Println(err)
		res := InvalidResponse{
			err.Error(),
		}

		b, _ := json.Marshal(res)
		w.Write(b)
		return

	}

})

var IndexHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Index page here mother fuckers"))
})
