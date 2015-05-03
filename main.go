package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cosmos-io/cosmos/dao"
	"github.com/cosmos-io/cosmos/router"
	"github.com/cosmos-io/cosmos/service"
	"github.com/cosmos-io/cosmos/worker"
	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
	"github.com/attilaolah/strict"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
)

var (
	cosmosPort  = getEnv("COSMOS_PORT", "8888")
	dbHost      = getEnv("INFLUXDB_HOST", "localhost")
	dbPort      = getEnv("INFLUXDB_PORT", "8086")
	dbUsername  = getEnv("INFLUXDB_USERNAME", "root")
	dbPassword  = getEnv("INFLUXDB_PASSWORD", "root")
	dbDatabase  = getEnv("INFLUXDB_DATABASE", "cosmos")
	dbShardConf = getEnv("INFLUXDB_SHARD_CONF", "./shard_config.json")
)

// to get an environment variable if it exists or default value
//
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

// to distingush text/html content type and others
// 
func contentType() martini.Handler {
	var (
		index []byte
		err   error
	)
	if martini.Env == "production" {
		index, err = ioutil.ReadFile("telescope/public/index.html")
	}

	if err != nil {
		fmt.Println("It is failed to read an index file.\n" + err.Error())
	}

	return func(r render.Render, req *http.Request, c martini.Context) {
		accept := strings.ToLower(req.Header.Get("Accept"))

		if strings.Contains(accept, "text/html") {
			if martini.Env == "development" {
				index, err = ioutil.ReadFile("telescope/public/index.html")
			}
			r.Header().Set(render.ContentType, "text/html; charset=utf-8")
			r.Data(http.StatusOK, index)
		}
	}
}

// to create an influxdb client
//
func createInfluxDBClient() *influxdbc.InfluxDB {
	dbc := influxdbc.NewInfluxDB(fmt.Sprintf("%s:%s", dbHost, dbPort), dbDatabase, dbUsername, dbPassword)
	file, err := ioutil.ReadFile(dbShardConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var conf influxdbc.ShardConfig
	err = json.Unmarshal(file, &conf)
	if err != nil {
		fmt.Println(err)
	}
	_, err = dbc.CreateDatabase(conf)
	if err != nil {
		fmt.Println("It is failed to create a database.\n" + err.Error())
	} else {
		fmt.Println("A database is created.")
	}

	return dbc
}

//
//
func cosmosService() martini.Handler {
	dbc := createInfluxDBClient()
	dao.Initialize(dbc)
	cosmos := service.NewCosmosService(5)

	newsFeedWorker := worker.NewNewsFeedWorker(cosmos.LifeTime, 30)
	newsFeedWorker.Run()

	return func(c martini.Context) {
		c.Map(cosmos)
		c.Next()
	}
}

// to run a martini app
//
func run() {
	m := martini.Classic()

	m.Handlers(
		gzip.All(),
		martini.Logger(),
		martini.Static("telescope/public"),
		strict.Strict,
		render.Renderer(),
		contentType(),
		cosmosService(),
	)

	m.Group("/v1", func(r martini.Router) {
		// get newsfeed
		r.Get("/newsfeeds",
			strict.Accept("application/json"),
			router.GetNewsFeeds)
		// get planet list
		r.Get("/planets",
			strict.Accept("application/json"),
			router.GetPlanets)

		r.Get("/planets/:planetName",
			strict.Accept("application/json"),
			router.GetPlanetMetrics)

		r.Get("/containers",
			strict.Accept("application/json"),
			router.GetContainers)

		// post container informations
		r.Post("/planets/:planetName/containers",
			strict.Accept("application/json"),
			strict.ContentType("application/json"),
			router.AddContainersOfPlanet)

		// get container list of planet
		r.Get("/planets/:planetName/containers",
			strict.Accept("application/json"),
			router.GetContainersOfPlanet)

		// get metrics of container
		r.Get("/planets/:planetName/containers/:containerName",
			strict.Accept("application/json"),
			router.GetContainerMetrics)
	})

	m.RunOnAddr(":" + cosmosPort)
}

func main() {
	run()
}
