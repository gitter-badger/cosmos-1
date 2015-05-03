package main

import (
	"os"
	"fmt"
    "path"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"

	"github.com/cosmos-io/cosmos/dao"
	"github.com/cosmos-io/cosmos/router"
	"github.com/cosmos-io/cosmos/service"
	"github.com/cosmos-io/cosmos/worker"
	"github.com/cosmos-io/influxdbc"
    
	"github.com/go-martini/martini"
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
func serveIndexHTML() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		accept := strings.ToLower(r.Header.Get("Accept"))
		if strings.Contains(accept, "text/html") {
            fp := path.Join("telescope", "public", "index.html")
            http.ServeFile(w, r, fp)
		}
	}
}

//
//
func serveCosmosService() martini.Handler {
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

// to run a martini app
//
func run() {
	m := martini.Classic()

	m.Handlers(
		martini.Logger(),
		martini.Static("telescope/public"),
		serveIndexHTML(),
		serveCosmosService(),
	)

	m.Group("/v1", func(r martini.Router) {
		// get newsfeed
		r.Get("/newsfeeds", router.GetNewsFeeds)

        // get planet list
		r.Get("/planets", router.GetPlanets)
        
		r.Get("/planets/:planetName", router.GetPlanetMetrics)
        
		r.Get("/containers", router.GetContainers)
        
		// get container list of planet
		r.Get("/planets/:planetName/containers", router.GetContainersOfPlanet)

		// get metrics of container
		r.Get("/planets/:planetName/containers/:containerName", router.GetContainerMetrics)

		// post container informations
		r.Post("/planets/:planetName/containers", router.AddContainersOfPlanet)
	})

	m.RunOnAddr(":" + cosmosPort)
}

func main() {
	run()
}
