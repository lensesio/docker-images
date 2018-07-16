FROM alpine
MAINTAINER Marios Andreopoulos <marios@landoop.com>

# Update, install tooling and some basic setup
RUN apk add --no-cache \
        bash \
        bash-completion \
        bzip2 \
        coreutils \
        curl \
        dumb-init \
        gettext \
        gzip \
        jq \
        libstdc++ \
        openjdk8-jre-base \
        openssl \
        supervisor \
        tar \
        wget \
    && echo "progress = dot:giga" | tee /etc/wgetrc \
    && mkdir -p /opt \
    && mkdir /etc/supervisord.d /etc/supervisord.templates.d

SHELL ["/bin/bash", "-c"]
WORKDIR /

# Install HortonWorks Schema Registry
ENV HW_SR_VERSION=0.5.2
ENV HW_SR_URL="https://github.com/hortonworks/registry/releases/download/v${HW_SR_VERSION}/hortonworks-registry-${HW_SR_VERSION}.tar.gz"
RUN wget "$HW_SR_URL" -O hw-registry.tar.gz \
    && mkdir -p /opt/hortonworks-registry \
    && tar xf /hw-registry.tar.gz --no-same-owner --strip-components=1 -C /opt/hortonworks-registry \
    && rm -f /hw-registry.tar.gz

ADD supervisord.conf /etc
ADD hw-registry.conf /etc/supervisord.d
ADD setup-and-run.sh /usr/local/bin
RUN chmod +x /usr/local/bin/setup-and-run.sh

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/local/bin/setup-and-run.sh"]
