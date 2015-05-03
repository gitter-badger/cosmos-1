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
    
    "github.com/gorilla/mux"
    "github.com/gorilla/context"
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

// A middleware to distingush text/html content type and others
// 
func serveIndexHTML(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := strings.ToLower(r.Header.Get("Accept"))
		if strings.Contains(accept, "text/html") {
            fp := path.Join("telescope", "public", "index.html")
            http.ServeFile(w, r, fp)
            return
		}
        next.ServeHTTP(w, r)
    })
}

// A middleware to serve cosmos instance
// 
func serveCosmos(next http.Handler) http.Handler {
	db := createInfluxDBClient()
	dao.Initialize(db)
	cosmos := service.NewCosmosService(5)

	newsFeedWorker := worker.NewNewsFeedWorker(cosmos.LifeTime, 30)
	newsFeedWorker.Run()

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        context.Set(r, "cosmos", cosmos)
        next.ServeHTTP(w, r)
    })
}

// to create an influxdb client
//
func createInfluxDBClient() *influxdbc.InfluxDB {
	db := influxdbc.NewInfluxDB(fmt.Sprintf("%s:%s", dbHost, dbPort), dbDatabase, dbUsername, dbPassword)
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
	_, err = db.CreateDatabase(conf)
	if err != nil {
		fmt.Println("It is failed to create a database.\n" + err.Error())
	} else {
		fmt.Println("A database is created.")
	}

	return db
}

// to run
//
func run() {
    staticFilesPath := path.Join("telescope", "public")
    
    mux := mux.NewRouter()
    mux.HandleFunc("/v1/newsfeeds", router.GetNewsFeeds).Methods("GET")
    mux.HandleFunc("/v1/planets/{planet}/containers/{container}", router.GetContainerMetrics).Methods("GET")
    mux.HandleFunc("/v1/planets/{planet}/containers", router.GetContainersOfPlanet).Methods("GET")
    mux.HandleFunc("/v1/planets/{planet}", router.GetPlanetMetrics).Methods("GET")
    mux.HandleFunc("/v1/planets", router.GetPlanets).Methods("GET")
    mux.HandleFunc("/v1/containers", router.GetContainers).Methods("GET")
    mux.HandleFunc("/v1/planets/{planet}/containers", router.AddContainersOfPlanet).Methods("POST")

    mux.PathPrefix("/").Handler(http.FileServer(http.Dir(staticFilesPath)))

    middlewares := serveIndexHTML(mux)
    middlewares = serveCosmos(middlewares)

    http.Handle("/", middlewares)
    http.ListenAndServe(":" + cosmosPort, nil)
}

func main() {
	run()
}
