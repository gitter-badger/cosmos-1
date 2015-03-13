FROM golang:1.4.2
MAINTAINER brann <brann@cosmos.io>

RUN mkdir -p /go/src/cosmos
WORKDIR /go/src/cosmos
ENV PATH /go/bin:$PATH

COPY ./cosmos_shard.conf /cosmos_shard.conf

ENV COSMOS_PORT 8080
ENV INFLUXDB_HOST influxdb
ENV INFLUXDB_PORT 8086
ENV INFLUXDB_USERNAME root
ENV INFLUXDB_PASSWORD root
ENV INFLUXDB_DATABASE cosmos
ENV INFLUXDB_SHARD_CONF /cosmos_shard.conf

EXPOSE 8080


# Install Godep
RUN go get github.com/tools/godep

COPY . /go/src/cosmos
RUN godep go install

CMD ["cosmos"]