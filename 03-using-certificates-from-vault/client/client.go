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
	// PathToCAFile is the path in the filesystem from where to read the root
	// certificate authority's certificate
	PathToCAFile = "/secret/chain.pem"
	// PathToCertFile is the path in the filesystem from where to read the
	// certificate file for 'time-client'
	PathToCertFile = "/secret/crt.pem"
	// PathToPrivateKeyFile is the path in the filesystem from where to read the
	// private key file for 'time-client'
	PathToPrivateKeyFile = "/secret/key.pem"
)

var (
	url       string
	tlsConfig *tls.Config
)

// getTime requrests the current time from 'time-server'
func getTime() (string, error) {
	// prepare a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// create a new http client using tlsConfig as its TLS configuration
	cli := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

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
	fs.StringVar(&url, "url", "https://time-server/", "url of the time server")
	fs.Parse(os.Args[1:])

	// read the certificate authority bundle from PathToCAFile
	caBundle, err := ioutil.ReadFile(PathToCAFile)
	if err != nil {
		log.Fatalf("failed to read ca bundle: %v", err)
	}

	// create a new x509 CertPool where to store the CA bundle
	// certificates in this bundle will be trusted as root CAs
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(caBundle); !ok {
		log.Fatalf("failed to read ca bundle: %v", err)
	}

	// load the certificate and private key for 'time-client' from
	// PathToCertFile and PathToPrivateKeyFile
	crt, err := tls.LoadX509KeyPair(PathToCertFile, PathToPrivateKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	// create a new tls.Config struct containing the desired TLS configuration:
	// - authenticate using the certificate read from PathToCertFile
	// - require TLSv1.2
	// - treat the certificates in the x509 CertPool created above as trusted
	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{crt},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      roots,
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
