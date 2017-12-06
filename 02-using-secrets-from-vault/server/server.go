package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// PathToUsername is the path in the filesystem from where to read the
	// username used to authenticate against 'time-server'
	PathToUsername = "/secret/username"
	// PathToPassword is the path in the filesystem from where to read the
	// password used to authenticate against 'time-server'.
	PathToPassword = "/secret/password"
)

var (
	port int

	username string
	password string
)

// getTime returns a string containing the current time in RFC3339 format
func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.IntVar(&port, "port", 80, "the port on which to listen")
	fs.Parse(os.Args[1:])

	// read the username from PathToUsername
	if val, err := ioutil.ReadFile(PathToUsername); err != nil {
		log.Fatalf("failed to read secret from %s: %v", PathToUsername, err)
	} else {
		username = string(val)
	}

	// read the password from PathToPassword
	if val, err := ioutil.ReadFile(PathToPassword); err != nil {
		log.Fatalf("failed to read secret from %s: %v", PathToPassword, err)
	} else {
		password = string(val)
	}

	// create a new http server that will listen on the specified port
	srv := http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
	}

	// add a handler to GET /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// retrieve basic authentication data
		user, pass, ok := r.BasicAuth()

		if ok && user == username && pass == password {
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
