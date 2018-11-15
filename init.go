/* HTTP BOOTSTRAP PACKGE
   This pakage define the entity and method that both Http server init and Middlewares will use
*/

package main

import (
	"net/http"
)

type Gate func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Gate) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
