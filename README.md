# h

[![Go Report Card](https://goreportcard.com/badge/github.com/josephspurrier/h)](https://goreportcard.com/report/github.com/josephspurrier/h)
[![GoDoc](https://godoc.org/github.com/josephspurrier/h?status.svg)](https://godoc.org/github.com/josephspurrier/h)

## Advanced Lightweight Go HTTP Handler Adapter

Simple to use.

Inspired by [mholt](https://github.com/mholt) and his project,
[caddy](https://github.com/mholt/caddy/wiki/Writing-a-Plugin:-HTTP-Middleware#writing-a-handler).

**h** provides a simple adapter to allow HTTP handler functions to return an int
HTTP status code and an error. This allows you to centralize the handling of
errors in one function.

Even [Andrew Gerrand](https://github.com/adg) uses a variation of this adapter
on the [The Go Blog](https://blog.golang.org/error-handling-and-go).