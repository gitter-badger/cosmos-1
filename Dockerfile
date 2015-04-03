FROM golang:1.4.2
MAINTAINER brann <brann@cosmos.io>

ENV PATH /go/bin:$PATH

COPY ./shard_config.json /shard_config.json

ENV COSMOS_PORT 8080
ENV INFLUXDB_HOST influxdb
ENV INFLUXDB_PORT 8086
ENV INFLUXDB_USERNAME root
ENV INFLUXDB_PASSWORD root
ENV INFLUXDB_DATABASE cosmos
ENV INFLUXDB_SHARD_CONF ./shard_config.json

EXPOSE 8080


# Install Godep
RUN go get github.com/tools/godep

# Copy source code
RUN mkdir -p /go/src/github.com/cosmos-io/cosmos
COPY . /go/src/github.com/cosmos-io/cosmos
WORKDIR /go/src/github.com/cosmos-io/cosmos

RUN godep go install

CMD ["cosmos"]