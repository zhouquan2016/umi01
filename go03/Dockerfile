FROM centos:centos8.4.2105
MAINTAINER zhqn
RUN mkdir /var/go03 /etc/go03
WORKDIR /var/go03
ENV GIN_MODE=release
ENV CONFIG=/etc/go03/
EXPOSE 8282
COPY go03 /var/go03/
COPY config.json /etc/go03/
COPY es.crt /etc/go03/
ENTRYPOINT ["/var/go03/go03", "-config=/etc/go03/"]