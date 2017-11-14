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
	url string

	username string
	password string
)

func getTime() (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(username, password)

	cli := http.Client{}

	res, err := cli.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

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
