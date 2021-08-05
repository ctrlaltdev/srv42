package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"srv42/utils"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const version = "1.2.0"

var (
	host    = flag.String("h", "127.0.0.1", "host to listen on")
	port    = flag.Int("p", 1337, "port to use")
	status  = flag.Int("s", 202, "status code to return")
	verbose = flag.Bool("v", false, "should srv42 print the full path and body")
	debug   = flag.Bool("d", false, "print debug information")

	wait = 60 * time.Second

	root *string
)

func ListenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(*status)
}

func Init() {

	args := flag.Args()
	utils.Debug(debug, fmt.Sprintf("Arguments | %v", args))

	if len(args) == 0 {

		utils.Debug(debug, "No arguments provided, fallback to listen only")
		return

	}

	utils.Debug(debug, fmt.Sprintf("Flags | host: %v - port: %v", *host, *port))

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

	srv := &http.Server{
		Handler:      l,
		Addr:         fmt.Sprintf("%s:%d", *host, *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Serving on: http://localhost:%d/\n", *port)

	go func() {
		err := srv.ListenAndServe()
		utils.CheckErr(err)
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

func main() {
	flag.Parse()

	Init()

	Serve()
}
