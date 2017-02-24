FROM centos:latest
MAINTAINER Marios Andreopoulos <marios@landoop.com>

RUN yum install -y epel-release \
    && yum install -y \
         gcc cpp glibc-devel glibc-headers kernel-headers libgomp libmpc mpfr \
         git \
         https://archive.landoop.com/third-party/librdkafka/librdkafka1-0.9.4.el7.centos.x86_64.rpm \
         https://archive.landoop.com/third-party/librdkafka/librdkafka-devel-0.9.4.el7.centos.x86_64.rpm

ENV DGOVERSION=1.8 DGOOS=linux DGOARCH=amd64
RUN curl https://storage.googleapis.com/golang/go${DGOVERSION}.${DGOOS}-${DGOARCH}.tar.gz \
    | tar xzf - -C /usr/local

ENV GOPATH=/opt/go PATH=/usr/local/go/bin:$PATH:/opt/go/bin
ADD . /opt/go/src/landoop/ais-generator

WORKDIR $GOPATH/src/landoop/ais-generator
RUN cd $GOPATH/src/landoop/ais-generator \
    && go get -v \
    && go build -v

ENV GODEBUG=cgocheck=0
ENTRYPOINT ["/opt/go/bin/ais-generator" ]
