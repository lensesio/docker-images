FROM centos:latest
MAINTAINER Marios Andreopoulos <marios@landoop.com>

RUN yum install -y \
        golang \
        git \
        https://archive.landoop.com/third-party/librdkafka/librdkafka1-0.9.4.el7.centos.x86_64.rpm \
        https://archive.landoop.com/third-party/librdkafka/librdkafka-devel-0.9.4.el7.centos.x86_64.rpm

ENV GOPATH=/opt/go PATH=$PATH:/opt/go/bin
ADD . /opt/go/src/landoop/ais-generator

WORKDIR $GOPATH/src/landoop/ais-generator
RUN cd $GOPATH/src/landoop/ais-generator \
    && go get -v \
    && go build -v

ENTRYPOINT ["/opt/go/bin/ais-generator" ]
