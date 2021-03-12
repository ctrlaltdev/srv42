package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

const (
	FormatPattern = "HTTP - %s - - [%s] \"%s %s %s\" %d %d \"%s\" %d\n"
)

var (
	port   = flag.Int("p", 1337, "port to use")
	status = flag.Int("s", 202, "status code to return")
)

func Handler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	w.WriteHeader(*status)

	fmt.Printf(FormatPattern,
		r.RemoteAddr,
		t.Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.URL.Path,
		r.Proto,
		*status,
		0,
		r.UserAgent(),
		time.Since(t),
	)
}

func Serve() {
	http.HandleFunc("/", Handler)

	fmt.Printf("Serving on port %d\n\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func main() {
	flag.Parse()

	Serve()
}
