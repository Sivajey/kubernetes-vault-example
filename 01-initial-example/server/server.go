package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// Username is the username used to authenticate against 'time-server'
	Username = "time-client"
	// Password is the password used to authenticate against 'time-server'
	Password = "safe#passw0rd!"
)

var (
	port int
)

// getTime returns a string containing the current time in RFC3339 format
func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.IntVar(&port, "port", 80, "the port on which to listen")
	fs.Parse(os.Args[1:])

	// create a new http server that will listen on the specified port
	srv := http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
	}

	// add a handler to GET /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// retrieve basic authentication data
		user, pass, ok := r.BasicAuth()

		if ok && user == Username && pass == Password {
			// if the username and password match respond with the current time
			log.Printf("got time request from '%s'", user)
			w.Write([]byte(getTime()))
		} else {
			// otherwise respond with 401
			w.Header().Set("WWW-Authenticate", `Basic realm="time-server"`)
			w.WriteHeader(401)
		}
	})

	// listen to and serve requests until an error occurs
	log.Fatal(srv.ListenAndServe())
}
