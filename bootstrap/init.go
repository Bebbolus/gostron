/* HTTP BOOTSTRAP PACKGE
    This pakage define the entity and method that both Http server init and Middlewares will use
 */

package bootstrap

import (
    "net/http"
)

type Gate func(http.HandlerFunc) http.HandlerFunc

// Chain applies Gate to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Gate) http.HandlerFunc {
    for _, m := range middlewares {
        f = m(f)
    }
    return f
}

