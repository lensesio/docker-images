A docker image that generates AIS data into a kafka topic.

A sample run could be:

    docker run --rm -it --net=host landoop/ais-generator \
        -messages 100000000 \
        -rate 25000 \
        -jitter 5000 \
        -bootstrap-servers broker1:9092,broker2:9092
        -schema-registry http://schema.registry:8081

The generator specifically produces _Class A Position Reports_. It includes
a real world recording of about 275000 such reports and rolls over once it
send them all.

The produce process follows a normal distribution with `mean=rate` and
`stddev=jitter`. Thus if you set `jitter=0` (default) you get a stable rate.

The reason we distribute it as a docker image is the librdkafka
bindings for Go, which make it non-trivial to build and use properly.
