package main

import (
	"net/http"
    "html/template"
)

type controller string

func (h controller) Fire(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("views/example.html") //setp 1
    t.Execute(w, "Programmer") //step 2
}

// Controller exported name
var Controller controller
