package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var HTTP_PORT int
var SECURE_TOKEN string
var HOSTNAME string

func init() {
	flag.IntVar(&HTTP_PORT, "p", 9001, "HTTP port")
	flag.StringVar(&SECURE_TOKEN, "token", "", "Secure token")
	flag.StringVar(&HOSTNAME, "hostname", "", "Secure token")
}

func main() {
	// Warn no auth
	if len(strings.TrimSpace(SECURE_TOKEN)) < 1 {
		log.Println("WARNING! Starting without secure token, no authentication enabled")
	}

	// Hostname
	if len(strings.TrimSpace(SECURE_TOKEN)) < 1 {
		HOSTNAME = getHostname()
		log.Printf("Detected hostname '%s'", HOSTNAME)
	}

	// Listen HTTP
	http.HandleFunc("/enqueue", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", HTTP_PORT), nil)
}

func getHostname() string {
	// Hostname
	name, err := os.Hostname()
	if err != nil {
		return ""
	}
	return name
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Read body
	bodyBytes, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		log.Println("%s", bodyErr)
		return
	}

	// Validate body
	bodyStr := string(bodyBytes)
	if len(bodyStr) < 1 {
		log.Println("Empty body")
		return
	}
}
