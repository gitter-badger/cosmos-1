package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/attilaolah/strict"
	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	db *influxdbc.InfluxDB

	cosmosPort  = getEnv("COSMOS_PORT", "8080")
	dbHost      = getEnv("INFLUXDB_HOST", "localhost")
	dbPort      = getEnv("INFLUXDB_PORT", "8086")
	dbUsername  = getEnv("INFLUXDB_USERNAME", "root")
	dbPassword  = getEnv("INFLUXDB_PASSWORD", "root")
	dbDatabase  = getEnv("INFLUXDB_DATABASE", "cosmos")
	dbShardConf = getEnv("INFLUXDB_SHARD_CONF", "./cosmos_shard.conf")
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
	return func(r render.Render, req *http.Request, c martini.Context) {
		accept := strings.ToLower(req.Header.Get("Accept"))
		if strings.Contains(accept, "text/html") {
			r.HTML(http.StatusOK, "index", nil)
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
			res := NewErrorJson(0, fmt.Sprintf("required parameter '%s' is missing", param))
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

func createInfluxDBConn() {
	db = influxdbc.NewInfluxDB(fmt.Sprintf("%s:%s", dbHost, dbPort), dbDatabase, dbUsername, dbPassword)
	file, err := ioutil.ReadFile(dbShardConf)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	var conf influxdbc.ShardSpace
	json.Unmarshal(file, &conf)
	_, err = db.CreateDatabase(conf)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func startServer() {
	m := martini.Classic()

	m.Handlers(
		martini.Logger(),
		martini.Static("web/public"),
		strict.Strict,
		render.Renderer(render.Options{
			Directory: "web/templates",
		}),
		contentTypeRouter(),
	)

	m.Group("/v1", func(r martini.Router) {
		r.Post("/:planet/containers",
			strict.Accept("application/json"),
			strict.ContentType("application/json"),
			requiredParams(),
			addContainers)
		r.Get("/:planet/containers",
			strict.Accept("application/json"),
			requiredParams("stime", "etime"),
			getContainers)
		r.Get("/planets",
			strict.Accept("application/json"),
			requiredParams(),
			getPlanets)
		r.Post("/planets",
			strict.Accept("application/json"),
			strict.ContentType("application/json"),
			requiredParams(),
			addPlanets)
	})

	if cosmosPort == "" {
		m.RunOnAddr(":3000")
	} else {
		m.RunOnAddr(":" + cosmosPort)
	}
}

func main() {
	createInfluxDBConn()
	startServer()
}
