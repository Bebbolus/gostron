package main

import (
	"net/http"
	"strings"
)

type middleware string

/*
	Method middleware ensures that url can only be requested with a specific method,
	else returns a 400 Bad Request
*/
func (m middleware) Pass(args string) func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			//MIDDLEWARE CORE THINGS
			acceptedMethods := strings.Split(args, "|")
			for _, v := range acceptedMethods {
				if r.Method == v {
					// Call the next middleware/handler in chain
					f(w, r)
					return
				}
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
}

// export as symbol named "Middleware"
var Middleware middleware
