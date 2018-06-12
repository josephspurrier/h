// Package h provides an adapter (F) and a settable function (ServeHTTP) to
// allow HTTP handler functions to return an int HTTP status code and an error.
package h

import (
	"net/http"
)

// F is an adapter to allow the use of a function that returns an int HTTP
// status code and an error as an HTTP handler.
type F func(http.ResponseWriter, *http.Request) (int, error)

// ServeHTTP is a settable function that receives the status and error from
// F.ServeHTTP.
var ServeHTTP = func(w http.ResponseWriter, r *http.Request, status int,
	err error) {
}

// ServeHTTP calls f(w, r) and passes the result to h.ServeHTTP.
func (fn F) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := fn(w, r)
	ServeHTTP(w, r, status, err)
}
