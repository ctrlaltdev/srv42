package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	FormatLog = "HTTP - %s - - [%s] \"%s %s %s\" %d %d \"%s\" %d\n"
)

var (
	port    = flag.Int("p", 1337, "port to use")
	status  = flag.Int("s", 202, "status code to return")
	verbose = flag.Bool("v", false, "should srv42 print the full path and body")
)

func Handler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	w.WriteHeader(*status)

	fmt.Printf(FormatLog,
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

	if *verbose {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n\n", string(requestDump))
	}
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
