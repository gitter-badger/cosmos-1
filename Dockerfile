FROM golang:1.4.2
MAINTAINER brann <brann@cosmos.io>

ENV PORT 8888
ENV INFLUXDB_HOST influxdb
ENV INFLUXDB_PORT 8086
ENV INFLUXDB_USERNAME root
ENV INFLUXDB_PASSWORD root
ENV INFLUXDB_DATABASE cosmos

EXPOSE 8888

COPY . /go/src/github.com/cosmos-io/cosmos
WORKDIR /go/src/github.com/cosmos-io/cosmos
RUN make

CMD ["./bin/cosmos"]