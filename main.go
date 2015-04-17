package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/attilaolah/strict"
	"github.com/cosmos-io/cosmos/router"
	"github.com/cosmos-io/cosmos/service"
	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
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

func getEnv(key, defVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defVal
	} else {
		return val
	}
}

func contentTypeRouter() martini.Handler {
	var (
		index []byte
		err   error
	)
	if martini.Env == "production" {
		index, err = ioutil.ReadFile("telescope/public/index.html")
	}

	if err != nil {
		fmt.Printf("\nFailed to read index file - %s\n\n", err.Error())
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

func requiredParams(params ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		valid := true
		var param string
		for _, param = range params {
			if query.Get(param) == "" {
				valid = false
				break
			}
		}
		if valid == false {
			res := map[string]interface{}{"error": fmt.Sprintf("required parameter '%s' is missing", param)}
			data, err := json.Marshal(res)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
		}
	}
}

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
		fmt.Println("Failed to create database - " + err.Error())
	} else {
		fmt.Println("Database has been created with sharding configuration")
	}

	return dbc
}

func cosmosService() martini.Handler {
	cosmos := service.NewCosmosService(createInfluxDBClient())

	return func(c martini.Context) {
		c.Map(cosmos)
		c.Next()
	}
}

func startServer() {
	m := martini.Classic()

	m.Handlers(
		martini.Logger(),
		martini.Static("telescope/public"),
		strict.Strict,
		render.Renderer(render.Options{
			Delims: render.Delims{Left: "{{%", Right: "%}}"},
		}),
		contentTypeRouter(),
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
			router.GetContainerInfo)
	})

	m.RunOnAddr(":" + cosmosPort)
}

func main() {
	startServer()
}
