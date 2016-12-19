# Fast Data KCQL service

A scala backend with `akka-http` and `reactive-kafka` capable to:

 * Stream messages from topics to clients via APIs (Rest, SSE and Web Sockets)
 * Execute KCQL statements
 * Allow time-travel

## Set up

Make sure you have access to

 * At least one Kafka Broker
 * At least one Zookeeper
 * At least on Schema Registry service

If you don't, just spin up our [`fast-data-dev`](https://github.com/landoop/fast-data-dev) docker

## Execute

```bash
docker run --rm -it -p 8080:8080 \
           -e brokers=broker1:9092,broker2:9092 \
           -e zookeeper=zk1:2181,zk2:2181/znode \
           -e registry=http://schema.registry.url:8081
           landoop/fast-data-backend
```

## Interact

### Web Sockets

Install the [`dark-websocket-terminal`](http://www.google.com/search?q=dark-websocket-terminal
) chrome plugin and execute:

    /connect ws://localhost:8080/api/kafka/ws?query=SELECT+*+FROM+_schemas
    /disconnect ws://localhost:8080/api/kafka/ws?query=SELECT+*+FROM+_schemas

    /connect ws://localhost:8080/api/kafka/ws/topicInfo
    /disconnect ws://localhost:8080/api/kafka/ws/topicInfo

### REST calls

    curl http://localhost:8080/api/kafka/rest/topicInfo

    curl http://localhost:8080/api/kafka/rest/kcql?query=SELECT+*+FROM+_schemas

## License

The project is licensed under the [BSL](http://www.landoop.com/bsl) license.
