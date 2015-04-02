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
	logDb *influxdbc.InfluxDB

	cosmosPort  = getEnv("COSMOS_PORT", "8080")
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

func createDBConn() {
	logDb = influxdbc.NewInfluxDB(fmt.Sprintf("%s:%s", dbHost, dbPort), dbDatabase, dbUsername, dbPassword)
	file, err := ioutil.ReadFile(dbShardConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("[database sharding configuration]")
	fmt.Println(string(file))

	var conf influxdbc.ShardConfig
	err = json.Unmarshal(file, &conf)
	if err != nil {
		fmt.Println(err)
	}
	_, err = logDb.CreateDatabase(conf)
	if err != nil {
		fmt.Println("Failed to create database")
		fmt.Println(err)
	} else {
		fmt.Println("database 'cosmos' is created with sharding configuration")
	}
}

func startServer() {
	m := martini.Classic()

	m.Handlers(
		martini.Logger(),
		martini.Static("telescope/public"),
		strict.Strict,
		render.Renderer(render.Options{
			Directory: "telescope/templates",
		}),
		contentTypeRouter(),
	)

	m.Group("/v1", func(r martini.Router) {
		// get planet list
		r.Get("/planets",
			strict.Accept("application/json"),
			getPlanets)

		// post container informations
		r.Post("/planets/:planet/containers",
			strict.Accept("application/json"),
			strict.ContentType("application/json"),
			addContainers)

		// get container list of planet
		r.Get("/planets/:planet/containers",
			strict.Accept("application/json"),
			//requiredParams("ttl"),
			getContainers)

		// get metrics of container
		r.Get("/planets/:planet/containers/:container",
			strict.Accept("application/json"),
			//requiredParams("interval"),
			getContainerInfo)
	})

	if cosmosPort == "" {
		m.RunOnAddr(":3000")
	} else {
		m.RunOnAddr(":" + cosmosPort)
	}
}

func main() {
	// obj := model.Container{Id: "ContainerId"}
	// obj.Ports = make([]*model.Port, 1)
	// privatePort := 80
	// publicPort := 80
	// //	portType := "TCP"

	// obj.Ports[0] = &model.Port{PrivatePort: &privatePort, PublicPort: &publicPort, Type: nil}

	// fmt.Println(obj)
	// t := reflect.TypeOf(obj)
	// v := reflect.ValueOf(obj)
	// tField := t.Field(0)
	// vField := v.Field(0)

	// fmt.Println(tField.Name)
	// fmt.Println(vField.Interface())

	// vals := MakeFieldPathAndValue(obj, ".")
	// for _, v := range vals {
	// 	fmt.Println(v)
	// }

	createDBConn()
	startServer()
}
