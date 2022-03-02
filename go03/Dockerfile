FROM centos:centos8.4.2105
MAINTAINER zhqn
RUN mkdir /var/godemo /etc/godemo
WORKDIR /var/godemo
ENV GIN_MODE=release
ENV GODEMO_CONFIG=/etc/godemo/
EXPOSE 8080
COPY go03 /var/godemo/
COPY config.json /var/godemo/
COPY es.crt /var/godemo/
ENTRYPOINT ["/var/godemo/go03"]