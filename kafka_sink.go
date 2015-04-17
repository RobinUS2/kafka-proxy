package main

// @author Robin Verlangen
// This will forward data to a Kafka cluster (queue)

import (
	"errors"
	"fmt"
	gokafka "github.com/Shopify/sarama"
	"log"
	"strings"
	"sync"
)

type KafkaSink struct {
	initMux   sync.RWMutex
	connected bool
	client    *gokafka.Client
	producer  gokafka.AsyncProducer
	buffer    chan string
	topic     string
	brokers   string
}

func (k *KafkaSink) Write(msg string) bool {
	select {
	case k.buffer <- msg:
		return true
	default:
		log.Println("Kafka sink buffer full, failed Write()")
		return false
	}
}

func (k *KafkaSink) _startWriter() bool {
	go func() {
		for {
			msg := <-k.buffer
			k._write(msg)
		}
	}()
	return true
}

func (k *KafkaSink) _write(msg string) bool {
	// Write to queue for buffering
	k.initMux.RLock()
	if k.connected == false {
		k.initMux.RUnlock()
		log.Println("Kafka not connected")
		return false
	}
	k.initMux.RUnlock()
	select {
	case k.producer.Input() <- &gokafka.ProducerMessage{Topic: k.topic, Key: nil, Value: gokafka.StringEncoder(msg)}:
		//log.Println("> message queued")
		return true
	case err := <-k.producer.Errors():
		log.Println(fmt.Sprintf("Failed to write into kafka: %s", err.Err))
		return false
	}
}

func (k *KafkaSink) Connect() (bool, error) {
	// Lock
	k.initMux.Lock()
	defer k.initMux.Unlock()

	// Once
	if k.connected {
		return false, errors.New("Kafka already connected")
	}

	// Connect
	log.Println("Kafka sink connecting")
	conf := gokafka.NewConfig()
	brokers := strings.Split(k.brokers, ",")
	conf.ClientID = fmt.Sprintf("kafka_sink_%s", HOSTNAME)
	client, clientErr := gokafka.NewClient(brokers, conf)
	if clientErr != nil {
		k.connected = false
		return false, clientErr
	}
	k.connected = true
	k.client = &client
	log.Println("Kafka sink connected")

	// Producer
	prodConf := gokafka.NewConfig()
	prodConf.Producer.Compression = gokafka.CompressionGZIP
	producer, producerErr := gokafka.NewAsyncProducer(brokers, prodConf)
	if producerErr != nil {
		return false, producerErr
	}
	k.producer = producer
	log.Println("Kafka producer initiated")

	return true, nil
}

func NewKafkaSink(topic string, brokers string) *KafkaSink {
	elm := &KafkaSink{
		topic:     topic,
		brokers:   brokers,
		connected: false,
		buffer:    make(chan string, 1000),
	}
	elm._startWriter()
	return elm
}
