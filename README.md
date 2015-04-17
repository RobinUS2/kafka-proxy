# kafka-proxy
Very simplistic HTTP proxy to forward data into kafka queues

Starting proxy
=====
One topic:
```
./run.sh --sinks my_topic=k001.cluster.com:9092,k002.cluster.com:9092
```

Multiple topics over multiple brokers/kafka-clusters:
```
./run.sh --sinks my_topic=k001.cluster.com:9092,k002.cluster.com:9092\;another_topic=kafka1.another-cluster.com:9092,kafka2.another-cluster.com:9092
```

Options
=====
The following options are available to start the proxy:
```
  -hostname="": Hostname override (default: automatic detection)
  -p=9001: HTTP port
  -sinks="": Configure sinks (mytopic=broker:port,broker:port;anothert_topic=broker:port)
  -token="": Authentication token
```

Forward data
=====
```
curl -XPOST --data "payload here" http://localhost:9001/enqueue?topic=my_topic
```

Example Use Cases
=====
- Forward data to a Kafka queue from a stateless programming language like PHP: this will increase throughput
- Centralize configuration and authentication
