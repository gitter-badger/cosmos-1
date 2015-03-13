About Cosmos
============
Display statistical informations of Docker container


Rest APIs
---------
Get all of the planet(=host) information in Cosmos
	[GET]  /v1/planets [Accept:application/json]

Get container metrics of planet from Cosmos
	[GET]  /v1/:planet/containers [Accept:application/json]

Post container metrics of planet to Cosmos
	[POST] /v1/:planet/containers [Content-type:application/json, Accept:application/json]
	Body : Raw
	[{
		"column" : "value",
		...
	}]