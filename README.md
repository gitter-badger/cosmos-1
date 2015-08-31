# Cosmos

[![Join the chat at https://gitter.im/cosmoshq/cosmos](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/cosmoshq/cosmos?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![Build Status](https://travis-ci.org/cosmoshq/cosmos.svg?branch=master)](https://travis-ci.org/cosmoshq/cosmos) [![GoDoc](https://godoc.org/github.com/cosmoshq/cosmos?status.svg)](https://godoc.org/github.com/cosmoshq/cosmos)

Cosmos is a container monitoring system. Cosmos can aggregate metrics of containers with [Curiosity](https://github.com/cosmoshq/curiosity). It also supports a modern dashboard.

## Quick start

You can run Cosmos simply.

```
$ docker run -d --name influxdb cosmosio/influxdb
$ docker run -d --link influxdb:influxdb -e INFLUXDB_HOST=influxdb --name cosmos cosmosio/cosmos:nightly
```

## Requirements
* InfluxDB (>= v0.9.0)
* Docker (>= v1.5.0)

## Debug

### InfluxDB

[InfluxDB](http://influxdb.com) is used in Cosmos. It is recommended to use [an InfluxDB container](https://registry.hub.docker.com/u/cosmoshq/influxdb/) with Cosmos. Of course, you can install InfluxDB in your local machine directly. If you do, please follow [the instruction](http://influxdb.com/download/).
```
$ docker run -p 8083:8083 -p 8086:8086 --expose 8090 --expose 8099 --rm --name influxdb cosmosio/influxdb
```

### Go

Cosmos is built with [Go](http://golang.org). [The latest version](https://golang.org/dl/) of Go is required when you debug.

* Go (>= 1.4.2)

```
$ git clone git@github.com:cosmoshq/cosmos
$ cd cosmos
$ make run
```

## Curiosity 

Curiosity is a container monitoring agent of Cosmos. You can run Curiosity simply with your COSMOS_HOST variable.

```
$ docker run -e COSMOS_HOST=127.0.0.1 --rm --name curiosity cosmosio/curiosity:nightly
```
See details: https://github.com/cosmoshq/curiosity

## Cosmos

<img src="https://raw.githubusercontent.com/cosmoshq/cosmos/master/screenshot.png">