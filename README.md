About Cosmos
============
Display statistical informations of Docker container


Rest APIs
---------
	[GET] /v1/:hostname/containers [Accept:application/json]
    
	[POST] /v1/:hostname/containers [Content-type:application/json, Accept:application/json]
	Body : Raw
	[{
		"column" : "value",
		...
	}]