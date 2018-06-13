# h

[![Go Report Card](https://goreportcard.com/badge/github.com/josephspurrier/h)](https://goreportcard.com/report/github.com/josephspurrier/h)
[![GoDoc](https://godoc.org/github.com/josephspurrier/h?status.svg)](https://godoc.org/github.com/josephspurrier/h)
[![Build Status](https://travis-ci.org/josephspurrier/h.svg)](https://travis-ci.org/josephspurrier/h)
[![Coverage Status](https://coveralls.io/repos/github/josephspurrier/h/badge.svg?branch=master)](https://coveralls.io/github/josephspurrier/h?branch=master)

## Advanced Lightweight Go HTTP Handler Adapter

**h** provides an adapter which allows HTTP handler functions to
return an int HTTP status code and an error. This technique allows you to
centralize the handling of errors via a customizable `ServeHTTP()` function.

Inspired by [mholt](https://github.com/mholt) and his project,
[caddy](https://github.com/mholt/caddy/wiki/Writing-a-Plugin:-HTTP-Middleware#writing-a-handler).

[Andrew Gerrand](https://github.com/adg) uses a custom HTTP handler similar to
this one (return just `error` instead of `int, error`) on the
[The Go Blog](https://blog.golang.org/error-handling-and-go).

## Usage

1. Import the package.

```go
import "github.com/josephspurrier/h"
```

2. Change each `http.HandleFunc()` to `http.Handle()` and wrap each HTTP handler
with `h.F()`.

```go
// Before
http.HandleFunc("/hello", Index)

// After
http.Handle("/hello", h.F(Index))
```

3. Add `(status int, err error)` as the return for each of your HTTP handler
functions. Also, modify each function to now return the proper values.

```go
// Before
func Index(w http.ResponseWriter, r *http.Request) {
	// ...
}

// After 
func Index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	// ...
}
```

4. Customize `h.ServeHTTP()` to fit your application needs.

```go
func main() {
	h.ServeHTTP = ServeHTTP
	http.Handle("/hello", h.F(Index))
	// ...
	http.ListenAndServe(":8080", nil)
}

// ServeHTTP handles all the HTTP handlers.
func ServeHTTP(w http.ResponseWriter, r *http.Request, status int, err error) {
	// Handle only errors.
	if status >= 400 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		r := new(randompackage.ErrResponse)
		r.Body.Status = http.StatusText(status)
		if err != nil {
			r.Body.Message = err.Error()
		}

		err := json.NewEncoder(w).Encode(r.Body)
		if err != nil {
			w.Write([]byte(`{"status":"Internal Server Error","message":"problem encoding JSON"}`))
			return
		}
	}

	// Only output 500 errors.
	if status >= 500 {
		if err != nil {
			log.Println(err)
		}
	}
}
```

## Full Example

Read through the comments on each function to see how it works.

You can build and run this application and then open your browser to any of the
following pages:

- http://localhost/
- http://localhost/created
- http://localhost/goodbye

```go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/josephspurrier/h"
)

// Hello will output "hello" with a 200 HTTP status code for requests to "/".
// For all other requests without a registered HTTP handler, it will show
// "page not found" with a 404 HTTP status code.
// Notice how the `return` is empty at the bottom of the function. Since the
// return values are named, they will return their zero values when an empty
// return is called. The default `h.ServeHTTP()` function also does not do
// anything with a status code less than 400 (400 and above are errors) so if
// you want it to be different than 200, you must use the `w.WriteHeader()` call
// prior to writing to the ResponseWriter.
func Hello(w http.ResponseWriter, r *http.Request) (status int, err error) {
	// Set the 404 error handler.
	if r.URL.Path != "/" {
		return http.StatusNotFound, errors.New("page not found")
	}
	fmt.Fprint(w, "hello")
	return
}

// Created will output "created" with a 201 HTTP status code.
// Notice how you have to specify `w.WriteHeader()` with the HTTP status code.
func Created(w http.ResponseWriter, r *http.Request) (status int, err error) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "created")
	return
}

// Goodbye will output "goodbye" with a 400 HTTP status code.
// Since you are returning the `http.StatusBadRequest` and the error, the
// `h.ServeHTTP()` function will determine how to handle the error. This allows
// you to specify how errors are handled in your application from one location.
func Goodbye(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return http.StatusBadRequest, errors.New("goodbye")
}

func main() {
	http.Handle("/", h.F(Hello))
	http.Handle("/created", h.F(Created))
	http.Handle("/goodbye", h.F(Goodbye))
	http.ListenAndServe(":8080", nil)
}
```

## In All Seriousness

You don't need to import this package into your project. It's actually designed
to show you how to use a custom HTTP handler. You can easily paste in this code
into your project and customize it to fit your needs.

```go
// F is an adapter to allow the use of a function that returns an int HTTP
// status code and an error as an HTTP handler.
type F func(http.ResponseWriter, *http.Request) (int, error)

// ServeHTTP calls f(w, r).
func (fn F) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := fn(w, r)
	if status >= 400 {
		if err != nil {
			http.Error(w, err.Error(), status)
		} else {
			http.Error(w, "", status)
		}
	}
}
```