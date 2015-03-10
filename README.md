About Cosmos
============
Display statistical informations of Docker container


Rest APIs
---------
	[GET] /v1/:hostname/containers
    
	[POST] /v1/:hostname/containers x-www-form-urlencoded
    parameter: data
    [{
        "column" : "value",
        ...
    }]