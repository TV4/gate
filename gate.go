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

import (
	"net/http"
)

type Gate struct {
	limitC chan struct{}
	waitC  chan struct{}
}

func New(limit, waitLimit int) *Gate {
	return &Gate{
		limitC: make(chan struct{}, limit),
		waitC:  make(chan struct{}, waitLimit),
	}
}

// BlockedHandler is the handler that is used when n < 1
var BlockedHandler = http.NotFoundHandler()

// Handler takes a http.Handler, limit number
func Handler(h http.Handler, limit int) http.Handler {
	return New(limit, 0).Handler(h)
}

// HandlerFunc takes a http.HandlerFunc and limit number
func HandlerFunc(h http.HandlerFunc, limit int) http.Handler {
	return Handler(h, limit)
}

func (g *Gate) Handler(h http.Handler) http.Handler {
	if cap(g.limitC) < 1 {
		return BlockedHandler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1
		if cap(g.waitC) > 0 {
			select {
			case g.waitC <- struct{}{}:
			default:
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			g.limitC <- struct{}{}
			<-g.waitC
		} else {
			g.limitC <- struct{}{}
		}
		defer func() { <-g.limitC }()

		h.ServeHTTP(w, r)

		// 2
		// select {
		// case g.limitC <- struct{}{}:
		// 	h.ServeHTTP(w, r)
		// 	<-g.limitC
		// case g.waitC <- struct{}{}:
		// 	h.ServeHTTP(w, r)
		// 	<-g.waitC
		// default:
		// 	http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		// }

		// 3
		// select {
		// case g.limitC <- struct{}{}:
		// 	defer func() { <-g.limitC }()
		// case g.waitC <- struct{}{}:
		// 	defer func() { <-g.waitC }()
		// default:
		// 	http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		// 	return
		// }

		// h.ServeHTTP(w, r)
	})
}

// Waiting returns the number of requests waiting at the gate
func (g *Gate) Waiting() int {
	return len(g.waitC)
}
