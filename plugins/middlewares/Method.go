package main

import (
	"net/http"
	"strings"
)

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Pass(m string) func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			//MIDDLEWARE CORE THINGS
			acceptedMethods := strings.Split(m, "|")
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
