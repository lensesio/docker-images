FROM landoop/fast-data-dev:latest
MAINTAINER Marios Andreopoulos <marios@landoop.com>

# Install ActiveMQ
ARG ACTIVEMQ_VERSION="5.14.2"
RUN wget "http://apache.mirrors.ovh.net/ftp.apache.org/dist/activemq/${ACTIVEMQ_VERSION}/apache-activemq-${ACTIVEMQ_VERSION}-bin.tar.gz" \
            -O /activemq.tgz \
    && mkdir -p /opt/activemq \
    && tar -xzf /activemq.tgz --no-same-owner --strip-components=1 -C /opt/activemq \
    && rm -rf /activemq.tgz

# Setup ActiveMQ
ARG ACTIVEMQ_ADMIN="admin: admin, admin"
ARG ACTIVEMQ_USER="user: user, user"
RUN sed -e '/admin:.*/d' \
        -e '/user:.*/d' \
        -i /opt/activemq/conf/jetty-realm.properties \
    && echo "$ACTIVEMQ_ADMIN" >> /opt/activemq/conf/jetty-realm.properties \
    && echo "$ACTIVEMQ_USER" >> /opt/activemq/conf/jetty-realm.properties

# Enable ActiveMQ (via supervisord)
RUN mkdir /etc/supervisord.d
ADD activemq.conf /etc/supervisord.d/
RUN echo "[include]" >> /etc/supervisord.conf \
    && echo 'files = /etc/supervisord.d/*.conf' >> /etc/supervisord.conf
ENV PATH $PATH:/opt/activemq/bin

# Install and configure AIS exporter
RUN apk add --no-cache go musl-dev git
ENV GOPATH=/opt/go PATH=$PATH:/opt/go/bin
ADD ais-mqtt /opt/go/src/landoop/ais-mqtt
RUN cd /opt/go/src/landoop/ais-mqtt/ \
    && go get -v \
    && go build -v

# Add MQTT connector configuration and Avro Schema
ADD mqtt-source.json classAPositionReportSchema.json /
