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
	url string

	username string
	password string
)

// getTime requrests the current time from 'time-server'
func getTime() (string, error) {
	// prepare a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// set the credentials for basic auth
	req.SetBasicAuth(username, password)

	// create a new http client
	cli := http.Client{}

	// send the request and grab the response
	res, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// if the request was not successfull, return an error
	if res.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	// read and return the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.StringVar(&url, "url", "http://time-server/", "url of the time server")
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

	// create a ticker that will be used to request time every five seconds
	t := time.NewTicker(5 * time.Second)

	for {
		<-t.C
		if time, err := getTime(); err != nil {
			log.Printf("failed to get time from server: %v", err)
		} else {
			log.Printf("got time from server: %s", time)
		}
	}
}
