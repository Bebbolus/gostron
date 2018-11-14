package main

import (
	"fmt"
	"net/http"
)

type handler string

func (h handler) Fire(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello FROM PLUGIN!!! ")

}

// export as symbol named "Handler"
var Handler handler
