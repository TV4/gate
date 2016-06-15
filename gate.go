/*

Package gate allows you to limit the number of busy handlers (at the gate)

Installation

Just go get the package:

    go get -u github.com/TV4/gate

Usage

A small usage example

    package main

    import (
    	"net/http"
    	"time"

    	"github.com/TV4/gate"
    )

    func index(w http.ResponseWriter, r *http.Request) {
    	time.Sleep(1 * time.Second)
    	w.Write([]byte("Hello, World!\n"))
    }

    func main() {
    	http.Handle("/", gate.HandlerFunc(index, 10))
    	http.ListenAndServe(":2016", nil)
    }

*/
package gate

import "net/http"

// BlockedHandler is the handler that is used when n < 1
var BlockedHandler = http.NotFoundHandler()

// Handler takes a http.Handler and limit number
func Handler(h http.Handler, n int) http.Handler {
	if n < 1 {
		return BlockedHandler
	}

	c := make(chan struct{}, n)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { <-c }()
		c <- struct{}{}
		h.ServeHTTP(w, r)
	})
}

// HandlerFunc takes a http.HandlerFunc and limit number
func HandlerFunc(h http.HandlerFunc, n int) http.Handler {
	return Handler(h, n)
}
