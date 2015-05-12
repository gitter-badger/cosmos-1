# Cosmos

[![Build Status](https://travis-ci.org/cosmos-io/cosmos.svg?branch=master)](https://travis-ci.org/cosmos-io/cosmos) [![GoDoc](https://godoc.org/github.com/cosmos-io/cosmos?status.svg)](https://godoc.org/github.com/cosmos-io/cosmos)

Cosmos is a container monitoring system. Cosmos can aggregate metrics of containers with [Curiosity](https://github.com/cosmos-io/curiosity). It also supports a modern dashboard.

## Quick start

You can run Cosmos simply.

```
$ docker run -d --name influxdb cosmosio/influxdb:0.9.0-rc29
$ docker run -d --link influxdb:influxdb -e INFLUXDB_HOST=influxdb --name cosmos cosmosio/cosmos:nightly
```

## Requirements
* InfluxDB (>= v0.9.0)
* Docker (>= v1.5.0)

## Debug

### InfluxDB

[InfluxDB](http://influxdb.com) is used in Cosmos. It is recommended to use [an InfluxDB container](https://registry.hub.docker.com/u/cosmosio/influxdb/) with Cosmos. Of course, you can install InfluxDB in your local machine directly. If you do, please follow [the instruction](http://influxdb.com/download/).
```
$ docker run -p 8083:8083 -p 8086:8086 --expose 8090 --expose 8099 --rm --name influxdb cosmosio/influxdb:0.9.0-rc29
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

## Cosmos

<img src="https://raw.githubusercontent.com/cosmos-io/cosmos/master/screenshot.png">