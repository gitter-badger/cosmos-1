# Cosmos

[![Build Status](https://travis-ci.org/cosmos-io/cosmos.svg?branch=master)](https://travis-ci.org/cosmos-io/cosmos) [![GoDoc](https://godoc.org/github.com/cosmos-io/cosmos?status.svg)](https://godoc.org/github.com/cosmos-io/cosmos)

Cosmos is a container monitoring system. Cosmos can aggregate metrics of containers with [Curiosity](https://github.com/cosmos-io/curiosity). It also supports a modern dashboard.

## Quick start

You can run Cosmos simply.

```
$ docker run -d --name influxdb cosmosio/influxdb
$ docker run -d --link influxdb:influxdb -e INFLUXDB_HOST=influxdb --name cosmos cosmosio/cosmos
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

## Curiosity 

Curiosity is a container monitoring agent of Cosmos. You can run Curiosity simply with your COSMOS_HOST variable.

```
$ docker run -e COSMOS_HOST=127.0.0.1 --rm --name curiosity cosmosio/curiosity:nightly
```
See details: https://github.com/cosmos-io/curiosity


## REST APIs

### Planets

Get all planets

```
GET /v1/planets 
Accept: application/json
```

### Containers

Get all containers

```
GET /v1/containers
Accept: application/json
```

Get all containers of planet

```
GET /v1/planet/:planet/containers
Accept: application/json
```

### Container

Get container

```
GET /v1/planet/:planetName/containers/:containerName
Accept: application/json
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
    "Names": ["cosmos"],
    "Ports": [{"PrivatePort": 80, "PublicPort": 80, "Type": "TCP" }],
    "Network": {"RxBytes": 8856494, "TxBytes: 102716"},
    "Cpu": {"TotalUtilization": 0.00, "PerCpuUtilization": [0.00,0.00,0.00,0.00]},
    "Memory": {"Limit": 2105901056, "Usage": 78176256}
  }

  ...
  
]
```