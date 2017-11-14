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
	Username = "time-client"
	Password = "safe#passw0rd!"
)

var (
	port int
)

func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.IntVar(&port, "port", 80, "the port on which to listen")
	fs.Parse(os.Args[1:])

	srv := http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if ok && user == Username && pass == Password {
			log.Printf("got time request from '%s'", user)
			w.Write([]byte(getTime()))
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="time-server"`)
			w.WriteHeader(401)
		}
	})

	log.Fatal(srv.ListenAndServe())
}
