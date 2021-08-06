package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"path"
	"srv42/utils"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const version = "1.3.1"

var (
	hostv6  = flag.String("hostv6", "[::1]", "IPv6 host to listen on")
	hostv4  = flag.String("hostv4", "127.0.0.1", "IPv4 host to listen on")
	port    = flag.Int("p", 1337, "port to use")
	status  = flag.Int("s", 202, "status code to return")
	verbose = flag.Bool("v", false, "should srv42 print the full path and body")
	debug   = flag.Bool("d", false, "print debug information")

	printVersion = flag.Bool("version", false, "print srv42 version")

	wait = 60 * time.Second

	root *string
)

func printv() {
	fmt.Printf("SRV42 - v%v\n", version)
	os.Exit(0)
}

func ListenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(*status)

	if *verbose {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n\n", string(requestDump))
	}
}

func Init() {

	args := flag.Args()
	utils.Debug(debug, fmt.Sprintf("Arguments | %v", args))

	if len(args) == 0 {

		utils.Debug(debug, "No arguments provided, fallback to listen only")
		return

	}

	utils.Debug(debug, fmt.Sprintf("Flags | hosts: %v + %v - port: %v", *hostv6, *hostv4, *port))

	if path.IsAbs(args[0]) {

		root = &args[0]
		utils.Debug(debug, fmt.Sprintf("Path is absolute: %s", *root))

	} else {

		cwd, err := os.Getwd()
		utils.CheckErr(err)

		absPath := path.Join(cwd, args[0])

		root = &absPath
		utils.Debug(debug, fmt.Sprintf("Path is not absolute | cwd: %s + %s : %s", cwd, args[0], *root))

	}
}

func Serve() {
	r := mux.NewRouter()

	if root != nil {

		fileServer := http.FileServer(http.Dir(*root))

		r.PathPrefix("/").Handler(fileServer).Methods(http.MethodGet, http.MethodOptions)

	} else {

		handler := http.HandlerFunc(ListenHandler)
		r.PathPrefix("/").Handler(handler).Methods(http.MethodGet, http.MethodOptions, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete)

	}

	r.Use(mux.CORSMethodMiddleware(r))

	l := handlers.LoggingHandler(os.Stdout, r)

	srvV6 := &http.Server{
		Handler:      l,
		Addr:         fmt.Sprintf("%s:%d", *hostv6, *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	srvV4 := &http.Server{
		Handler:      l,
		Addr:         fmt.Sprintf("%s:%d", *hostv4, *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Serving on: http://localhost:%d/\n", *port)

	go func() {
		err := srvV6.ListenAndServe()
		utils.CheckErr(err)
	}()

	go func() {
		err := srvV4.ListenAndServe()
		utils.CheckErr(err)
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srvV6.Shutdown(ctx)
	srvV4.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

func main() {
	flag.Parse()

	if *printVersion {
		printv()
	}

	Init()

	Serve()
}
