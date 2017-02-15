A docker image that generates AIS data into a kafka topic.

A sample run could be:

    docker run --rm -it --net=host landoop/ais-generator \
        -messages 100000000 \
        -rate 25000 \
        -bootstrap-servers broker1:9092,broker2:9092
        -schema-registry http://schema.registry:8081
