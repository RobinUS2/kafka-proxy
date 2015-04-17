package main

import (
	"encoding/json"
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
var SINKS string
var sinkMap map[string]*KafkaSink

func init() {
	flag.IntVar(&HTTP_PORT, "p", 9001, "HTTP port")
	flag.StringVar(&SECURE_TOKEN, "token", "", "Secure token")
	flag.StringVar(&HOSTNAME, "hostname", "", "Secure token")
	flag.StringVar(&SINKS, "sinks", "", "Configure sinks (mytopic=broker:port,broker:port;anothert_topic=broker:port)")
	flag.Parse()
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

	// Connect sinks
	sinkMap = make(map[string]*KafkaSink)
	topics := strings.Split(SINKS, ";")
	for _, topicDef := range topics {
		split := strings.Split(topicDef, "=")
		if len(split) != 2 {
			log.Println("WARNING! Sink definition format invalid")
			continue
		}
		var topic string = split[0]
		var brokers string = split[1]
		if len(strings.TrimSpace(topic)) < 1 {
			log.Println("WARNING! Empty topic in configuration")
			continue
		}
		if len(strings.TrimSpace(brokers)) < 1 {
			log.Println("WARNING! Empty broker list in configuration")
			continue
		}

		// Connect
		sinkMap[topic] = NewKafkaSink(topic, brokers)
		sinkMap[topic].Connect()
	}

	// Listen HTTP
	http.HandleFunc("/enqueue", handler)
	http.HandleFunc("/stats", statsHandler)
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

func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Token
	if r.URL.Query().Get("token") != SECURE_TOKEN {
		w.WriteHeader(401) // Not authorized
		return
	}

	// Stats
	stats := make(map[string]map[string]interface{})
	for topic, sink := range sinkMap {
		stats[topic] = make(map[string]interface{})
		stats[topic]["messages"] = sink.messageCount
	}
	bytes, _ := json.Marshal(stats)
	fmt.Fprintf(w, "%s", bytes)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Token
	if r.URL.Query().Get("token") != SECURE_TOKEN {
		w.WriteHeader(401) // Not authorized
		return
	}

	// Read body
	bodyBytes, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		w.WriteHeader(400) // Bad request
		log.Println("%s", bodyErr)
		return
	}

	// Validate body
	bodyStr := string(bodyBytes)
	if len(bodyStr) < 1 {
		w.WriteHeader(400) // Bad request
		log.Println("Empty body")
		return
	}

	// Topic
	reqTopic := strings.TrimSpace(r.URL.Query().Get("topic"))
	if len(reqTopic) < 1 || sinkMap[reqTopic] == nil {
		w.WriteHeader(400) // Bad request
		log.Println("Topic not found")
		return
	}

	// Insert
	sinkMap[reqTopic].Write(bodyStr)
}
