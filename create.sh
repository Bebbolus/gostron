#!/bin/bash

create_handler(){
	cat  <<EOF > plugins/controllers/$1.go
package main

import (
	"fmt"
	"net/http"
)

type controller string

func (h controller) Fire(w http.ResponseWriter, r *http.Request) {
	## YOUR HANDLER CODE GOES HERE
}

// Controller exported name
var Controller controller
EOF
}

create_middleware(){
	cat << EOF > plugins/middlewares/$1.go 
package main

import (
	"net/http"
	"strings"
)

type middleware string

func (m middleware) Pass(args string) func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			//MIDDLEWARE CORE THINGS
			if 1 ==1 {
				// Call the next middleware/handler in chain
				f(w, r)
				return
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
}

// Middleware exported symbol
var Middleware middleware

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

