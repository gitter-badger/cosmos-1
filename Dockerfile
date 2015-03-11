FROM golang:latest
MAINTAINER brann <brann@cosmos.io>

# Install Cosmos
ENV COSMOS_PORT 8080
ENV INFLUXDB_HOST influxdb
ENV INFLUXDB_PORT 8086
ENV INFLUXDB_USERNAME root
ENV INFLUXDB_PASSWORD root
ENV INFLUXDB_DATABASE cosmos

RUN go get github.com/cosmos-io/cosmos && go install github.com/cosmos-io/cosmos

EXPOSE 8081
CMD ["/go/bin/cosmos"]



#VOLUME ["/data"]
#CMD ["/run.sh"]
