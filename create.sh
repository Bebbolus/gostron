#!/bin/bash

create_handler(){
	cat  <<EOF > plugins/handlers/$1.go
package main

import (
	"net/http"
)

type handler string

func (h handler) Fire(w http.ResponseWriter, r *http.Request) {
	## YOUR HANDLER CODE GOES HERE
}

// export as symbol named "Handler"
var Handler handler
EOF
}

create_middleware(){
	cat << EOF > plugins/middlewares/$1.go 
package main

import (
	bootstrap "github.com/Bebbolus/gostron/bootstrap"
	"net/http"
	"strings"
)

func Pass(m string) bootstrap.Gate {
	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			//MIDDLEWARE CORE THINGS
			if 1==1{
				// Call the next middleware/handler in chain
				f(w, r)
				return
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
}
EOF
}

case $1 in
    handler)
        create_handler $2
        exit 1
        ;;
    middleware)
        create_middleware $2
        exit 1
        ;;
    *)
        echo "Specify a right target"
        exit 1
        ;;
esac

