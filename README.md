# Cosmos
You are a DevOps engineer and also a pioneer to a Cosmos. Cosmos is a monitoring system for containers. Cosmos can aggregate metrics and logs from a drone called Curiosity to an InfluxDB and analyze the data. You can find out more about [Curiosity](https://github.com/cosmos-io/curiosity).

## Run
Docker is the best option to run Cosmos.
```
$ docker run --rm --name influxdb cosmosio/influxdb
$ docker run --link influxdb:influxdb -e INFLUXDB_HOST=influxdb --rm --name cosmos cosmosio/cosmos
$ docker run --link cosmos:cosmos -e COSMOS_HOST=cosmos -e COSMOS_PLANET_NAME=Mars --rm --name curiosity cosmosio/curiosity
```

Rest APIs
---------
Get all of the planet(host) information in Cosmos
	[GET]  /v1/planets [Accept:application/json]

Post planet(host) information to Cosmos
     	[POST] /v1/planets [Content-type:application/json, Accept:application/json]
	Body : Raw
	{
		"column" : "value",
		...
	}

Get container metrics of planet from Cosmos
	[GET]  /v1/:planet/containers [Accept:application/json]
	Requried Parameters
	stime : start time of metrics in seconds
	etime : end time of metrics in seconds

Post container metrics of planet to Cosmos
	[POST] /v1/:planet/containers [Content-type:application/json, Accept:application/json]
	Body : Raw
	[{
		"column" : "value",
		...
	}]
