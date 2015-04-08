# Cosmos

[![Build Status](https://travis-ci.org/cosmos-io/cosmos.svg?branch=master)](https://travis-ci.org/cosmos-io/cosmos)

You are a DevOps engineer and also a pioneer to a Cosmos. Cosmos is a monitoring system for containers. Cosmos can aggregate metrics and logs from a drone called Curiosity to an InfluxDB and analyze the data. You can find out more about [Curiosity](https://github.com/cosmos-io/curiosity).

## Run
Docker is the best option to run Cosmos.
```
# Start influxdb container used by Cosmos
# See details : https://github.com/cosmos-io/influxdb-dockerfile
$ docker run --rm --name influxdb cosmosio/influxdb

# Start cosmos container 
$ docker run --link influxdb:influxdb -e INFLUXDB_HOST=influxdb --rm --name cosmos cosmosio/cosmos

# Start curiosity container for collecting metrics of containers on current Host(Planet)
# See details : https://github.com/cosmos-io/curiosity
# If you use same planet for cosmos and curiosity, then
$ docker run --link cosmos:cosmos -v /var/run/docker.sock:/var/run/docker.sock -e COSMOS_HOST=cosmos -e COSMOS_PLANET_NAME=Mars --rm --name curiosity cosmosio/curiosity
# Otherwise,
$ docker run -v /var/run/docker.sock:/var/run/docker.sock -e COSMOS_HOST=some.where --rm --name curiosity cosmosio/curiosity
```

Rest APIs
---------
Get all of the planet(host) information
	[GET]  /v1/planets [Accept:application/json]

Get all of the container metrics
    	[GET] /v1/containers [Accept:application/json]

Get container list of planet
	[GET]  /v1/planet/:planet/containers [Accept:application/json]

Post container metrics of planet
	[POST] /v1/planet/:planet/containers [Content-Type:application/json, Accept:application/json]

	Body :
	[{
	    "Id": "52236af62ef96f960611b6a9276d3be9800a9f04a497d24a4a7dfb1c74b23be6",
	    "Image": "cosmosio/rust:nightly",
	    "Status": "Up 2 minutes",
	    "Command": "bash",
	    "Created": 1428047585,
	    "Names": ["/cosmos"],
	    "Ports": [{
		"PrivatePort": 80,
		"PublicPort": 80,
		"Type": "TCP"
	    }],
	    "Stats": {
		"Network": {
		   "RxBytes":8856494,
		   "TxBytes":102716
		},
		"Cpu": {
		   "TotalUtilization": 11.2,
		   "PerCpuUtilization": [5.2, 6.0]
		},
		"Memory": {
		   "Limit":2105901056,
		   "Usage":78176256
	        }
	    }
	}]

Get metrics of container
	[GET] /v1/planet/:planet/containers/:container_name [Accept:application/json]


<img src="https://raw.githubusercontent.com/cosmos-io/cosmos/master/screenshot.png">