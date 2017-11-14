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
	Username = "time-client"
	Password = "safe#passw0rd!"
)

var (
	url string
)

func getTime() (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(Username, Password)

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
