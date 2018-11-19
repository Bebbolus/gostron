package main

import (
	"fmt"
	"net/http"
)

type controller string

func (h controller) Fire(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello FROM CONTROLLER PLUGIN!!!")

}

// Controller exported name
var Controller controller
