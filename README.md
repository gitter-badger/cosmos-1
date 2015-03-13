# Cosmos
You are a DevOps engineer and also a pioneer to a Cosmos. Cosmos is a monitoring system for containers. Cosmos can aggregate metrics and logs from a drone called Curiosity to an InfluxDB and analyze the data. You can find out more about [Curiosity](https://github.com/cosmos-io/curiosity).

## Run
Docker is the best option to run Cosmos.
```
$ docker run --rm --name influxdb cosmosio/influxdb
$ docker run --link influxdb:influxdb -e INFLUXDB_HOST=influxdb --rm --name cosmos cosmosio/cosmos
$ docker run --link cosmos:cosmos -e COSMOS_HOST=cosmos -e COSMOS_PLANET_NAME=Mars --rm --name curiosity cosmosio/curiosity
```

## Rest APIs
	[GET] /v1/:hostname/containers [Accept:application/json]
    
	[POST] /v1/:hostname/containers [Content-type:application/json, Accept:application/json]
	Body : Raw
	[{
		"column" : "value",
		...
	}]