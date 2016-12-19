FROM alpine:latest
MAINTAINER Marios Andreopoulos <marios@landoop.com>

WORKDIR /

RUN apk add --no-cache openjdk8-jre-base

ENV SBT_OPTS="-Xms512M -Xmx1536M -Xss1M -XX:+CMSClassUnloadingEnabled"
RUN apk add --no-cache git wget \
    && git clone https://github.com/Landoop/fast-data-backend.git /fast-data-backend \
    &&  wget https://repo.typesafe.com/typesafe/ivy-releases/org.scala-sbt/sbt-launch/0.13.13/sbt-launch.jar \
    && cd /fast-data-backend \
    && java $SBT_OPTS -jar /sbt-launch.jar assembly \
    && mv /fast-data-backend/target/scala-2.11/fast.data.kcql.service-assembly-1.0.jar / \
    && rm -rf /sbt-launch.jar /root/.sbt /root/.ivy2 rm -rf /fast-data-backend \
    && apk --no-cache del git wget

USER nobody
EXPOSE 8080
ENV brokers=localhost:9092 zookeeper=localhost:2181 registry=http://localhost:8081 port=8080

ENTRYPOINT /usr/bin/java \
             -Dkafka.brokers="$brokers" \
             -Dkafka.zookeeper="$zookeeper" \
             -Dkafka.schema-registry-url="$registry" \
             -Dport="$port" \
             -jar /fast.data.kcql.service-assembly-1.0.jar
