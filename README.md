# Cosmos

[![Build Status](https://travis-ci.org/cosmos-io/cosmos.svg?branch=master)](https://travis-ci.org/cosmos-io/cosmos)

Cosmos is a container monitoring system. Cosmos can aggregate metrics of containers with [Curiosity](https://github.com/cosmos-io/curiosity). It also supports a modern dashboard.

## Quick start

You can run Cosmos simply.

```
$ docker run -d --rm --name influxdb cosmosio/influxdb
$ docker run -d --link influxdb:influxdb -e INFLUXDB_HOST=influxdb --rm --name cosmos cosmosio/cosmos
```

## Requirements
* InfluxDB (>= v0.8.8)
* Docker (>= v1.5.0)

## Debug

### InfluxDB

[InfluxDB](http://influxdb.com) is used in Cosmos. It is recommended to use [a InfluxDB container](https://registry.hub.docker.com/u/cosmosio/influxdb/) with Cosmos. Of course, you can install InfluxDB in your local machine directly. If you do, please follow [the instruction](http://influxdb.com/download/).
```
$ docker run -p 8083:8083 -p 8086:8086 --expose 8090 --expose 8099 --rm --name influxdb cosmosio/influxdb
```

### Go

Cosmos is built with [Go](http://golang.org). [The latest version](https://golang.org/dl/) of Go is required when you debug.

* Go (>= 1.4.2)

```
$ go get github.com/cosmos-io/cosmos
```

## Run

```
$ docker run -p 8083:8083 -p 8086:8086 --expose 8090 --expose 8099 --rm --name influxdb cosmosio/influxdb
$ docker run --link influxdb:influxdb -e INFLUXDB_HOST=influxdb --rm --name cosmos cosmosio/cosmos
```

## Curiosity

Curiosity is a container monitoring agent of Cosmos. You can run Curiosity simply with your COSMOS_HOST variable.

```
$ docker run -e COSMOS_HOST=127.0.0.1 --rm --name curiosity cosmosio/curiosity:nightly
```

## REST APIs

### Planets

```
GET /v1/planets HTTP/1.1
Accept: application/json
```

### Containers

* ttl : time to live for containers (ex, 30m, 1h, 7d)

```
GET /v1/planet/:planet/containers
Accept: application/json

{"ttl":"30m"}
```

### Container

* interval : time interval from now (ex, 30m, 1h, 7d)

```
GET /v1/planet/:planet/containers/:container
Accept: application/json

{"interval":"30m"}
```

### Create containers

```
POST /v1/planet/:planet/containers
Accept: application/json
Content-Type: application/json

[
  {
    "Id": "52236af62ef96f960611b6a9276d3be9800a9f04a497d24a4a7dfb1c74b23be6",
    "Image": "cosmosio/cosmos:latest",
    "Status": "Up 2 minutes",
    "Command": "bash",
    "Created": "1428047585",
    "Names": ["/cosmos"],
    "Ports": [],
    "Network": {"RxBytes": 8856494, "TxBytes: 102716"},
    "Cpu": {"TotalUtilization": 0.00, "PerCpuUtilization": [0.00,0.00,0.00,0.00]},
    "Memory": {"Limit": 2105901056, "Usage": 78176256}
  }

  ...
  
]
```

## Screenshot

<img src="https://raw.githubusercontent.com/cosmos-io/cosmos/master/screenshot.png">