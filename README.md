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

[InfluxDB](http://influxdb.com) is also used in Cosmos. It is recommended to use [a InfluxDB container](https://registry.hub.docker.com/u/cosmosio/influxdb/) with Cosmos. Of course, you can install InfluxDB in your local machine directly. If you do, please follow [the instruction](http://influxdb.com/download/).
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
Get all of the planet(host) information
	[GET]  /v1/planets [Accept:application/json]


Get container list of planet
	[GET]  /v1/planet/:planet/containers [Accept:application/json]
	Required Parameters
		- ttl : time to live for containers (ex, 30m, 1h, 7d)

Post container metrics of planet
	[POST] /v1/planet/:planet/containers [Content-Type:application/json, Accept:application/json]
	Body : Raw
	{ 

	  "nginx" : {
				"column" : "value",
				...
	  },

	  "mysql" : {
	  			"column" : "value",
	  			...
	  }, 
	  ...	  
	}

Get metrics of container
	[GET] /v1/planet/:planet/containers/:container_name [Accept:application/json]
	Required parameters
		- interval : time interval from now (ex, 30m, 1h, 7d)

<img src="https://raw.githubusercontent.com/cosmos-io/cosmos/master/screenshot.png">