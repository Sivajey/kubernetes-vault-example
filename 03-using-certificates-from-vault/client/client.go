package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	PathToCAFile         = "/secret/chain.pem"
	PathToCertFile       = "/secret/crt.pem"
	PathToPrivateKeyFile = "/secret/key.pem"
)

var (
	url       string
	tlsConfig *tls.Config
)

func getTime() (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	cli := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

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
	fs.StringVar(&url, "url", "https://time-server/", "url of the time server")
	fs.Parse(os.Args[1:])

	caBundle, err := ioutil.ReadFile(PathToCAFile)
	if err != nil {
		log.Fatalf("failed to read ca bundle: %v", err)
	}

	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(caBundle); !ok {
		log.Fatalf("failed to read ca bundle: %v", err)
	}

	crt, err := tls.LoadX509KeyPair(PathToCertFile, PathToPrivateKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{crt},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      roots,
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
