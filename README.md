# kafka-proxy
Very simplistic HTTP proxy to forward data into kafka queues

Starting
=====
One topic:
```
./run.sh --sinks my_topic=k001.cluster.com:9092,k002.cluster.com:9092
```

Multiple topics over multiple brokers/kafka-clusters:
```
./run.sh --sinks my_topic=k001.cluster.com:9092,k002.cluster.com:9092\;another_topic=kafka1.another-cluster.com:9092,kafka2.another-cluster.com:9092
```
