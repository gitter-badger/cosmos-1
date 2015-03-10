FROM golang:latest
MAINTAINER brann <brann@cosmos.io>

# Install Cosmos
ENV COSMOS_PORT 8080
ENV INFLUXDB_HOST localhost
ENV INFLUXDB_PORT 8086
ENV INFLUXDB_USERNAME root
ENV INFLUXDB_PASSWORD root
ENV INFLUXDB_DATABASE cosmos

RUN go get github.com/cosmos-io/cosmos && go install github.com/cosmos-io/cosmos
CMD ["$GOPATH/bin/cosmos"]


# Raft port (for clustering, don't expose publicly!)
#EXPOSE 8090

# Protobuf port (for clustering, don't expose publicly!)
#EXPOSE 8099

#VOLUME ["/data"]
#CMD ["/run.sh"]
