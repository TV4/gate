# gate

[![Build Status](https://travis-ci.org/TV4/gate.svg?branch=master)](https://travis-ci.org/TV4/gate)
[![Go Report Card](https://goreportcard.com/badge/github.com/TV4/gate)](https://goreportcard.com/report/github.com/TV4/gate)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/TV4/gate)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/TV4/gate#license-mit)

Limit the number of busy handlers (at the gate)

## Installation

    go get -u github.com/TV4/gate

## Usage

```go
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
```

```bash
$ wrk http://localhost:2016/
Running 10s test @ http://localhost:2016/
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.00s     1.41ms   1.01s    70.00%
    Req/Sec     4.20      0.41     5.00     80.00%
  100 requests in 10.07s, 12.79KB read
Requests/sec:      9.93
Transfer/sec:      1.27KB
```

## License (MIT)

Copyright (c) 2016-2017 TV4

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:

> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
