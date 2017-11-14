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
	PathToUsername = "/secret/username"
	PathToPassword = "/secret/password"
)

var (
	port int

	username string
	password string
)

func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.IntVar(&port, "port", 80, "the port on which to listen")
	fs.Parse(os.Args[1:])

	if val, err := ioutil.ReadFile(PathToUsername); err != nil {
		log.Fatalf("failed to read secret from %s: %v", PathToUsername, err)
	} else {
		username = string(val)
	}

	if val, err := ioutil.ReadFile(PathToPassword); err != nil {
		log.Fatalf("failed to read secret from %s: %v", PathToPassword, err)
	} else {
		password = string(val)
	}

	srv := http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if ok && user == username && pass == password {
			log.Printf("got time request from '%s'", user)
			w.Write([]byte(getTime()))
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="time-server"`)
			w.WriteHeader(401)
		}
	})

	log.Fatal(srv.ListenAndServe())
}
