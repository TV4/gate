package gate

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testCases = []struct {
	body string
	want string
	code int
	reqs int
	n    int
	min  time.Duration
}{
	{"Hello, World!", "404 page not found\n", 404, 1, 0, 0},
	{"Hello, World!", "Hello, World!", 200, 1, 2, 50 * time.Millisecond},
	{"Hello, World!", "Hello, World!", 200, 2, 1, 50 * time.Millisecond},
	{"Hello, World!", "Hello, World!", 200, 2, 2, 100 * time.Millisecond},
	{"Hello, World!", "Hello, World!", 200, 2, 3, 100 * time.Millisecond},
	{"Hello, World!", "Hello, World!", 200, 3, 1, 150 * time.Millisecond},
	{"Hello, World!", "Hello, World!", 200, 4, 1, 200 * time.Millisecond},
}

func TestHandler(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	gh := Handler(h, 1)

	ts := httptest.NewServer(gh)
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if got, want := r.StatusCode, http.StatusTeapot; got != want {
		t.Fatalf("r.StatusCode = %d, want %d", got, want)
	}
}

func TestHandlerFunc(t *testing.T) {
	for i, tt := range testCases {
		h := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tt.body))
		}

		gh := HandlerFunc(h, tt.n)

		r, _ := http.NewRequest("GET", "", nil)
		w := httptest.NewRecorder()

		gh.ServeHTTP(w, r)

		if got, want := w.Code, tt.code; got != want {
			t.Fatalf("[%d] w.Code = %d, want %d", i, got, want)
		}

		if got, want := w.Body.String(), tt.want; got != want {
			t.Fatalf("[%d] w.Body.String = %q, want %q", i, got, want)
		}
	}
}

func TestGate_Handler(t *testing.T) {
	for i, tt := range testCases {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.Write([]byte(tt.body))
		})

		gh := New(tt.n).Handler(h)

		ts := httptest.NewServer(gh)
		defer ts.Close()

		start := time.Now()

		for ri := 0; ri < tt.reqs; ri++ {
			r, err := http.Get(ts.URL)
			if err != nil {
				t.Fatalf("[%d:%d] unexpected err: %v", i, ri, err)
			}
			defer r.Body.Close()

			if got, want := r.StatusCode, tt.code; got != want {
				t.Fatalf("[%d:%d] r.StatusCode = %d, want %d", i, ri, got, want)
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("[%d:%d] unexpected err: %v", i, ri, err)
			}

			if got, want := string(body), tt.want; got != want {
				t.Fatalf("[%d:%d] string(body) = %q, want %q", i, ri, got, want)
			}
		}

		elapsed := time.Since(start)

		if elapsed < tt.min {
			t.Errorf("[%d] elapsed: %v, want %v", i, elapsed, tt.min)
		}
	}
}
