FROM centos:latest
MAINTAINER Marios Andreopoulos <marios@landoop.com>

ENV BORG_RELEASE=1.0.8

ADD https://github.com/borgbackup/borg/releases/download/${BORG_RELEASE}/borg-linux64 /usr/bin/borg
RUN chmod +x /usr/bin/borg

ENTRYPOINT ["/usr/bin/borg"]
