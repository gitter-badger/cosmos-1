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