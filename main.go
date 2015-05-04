package main

import (
	"os"
	"fmt"
    "path"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"

    "github.com/cosmos-io/cosmos/context"
	"github.com/cosmos-io/cosmos/dao"
	"github.com/cosmos-io/cosmos/router"
	"github.com/cosmos-io/cosmos/service"
	"github.com/cosmos-io/cosmos/worker"
	"github.com/cosmos-io/influxdbc"
    
    "github.com/gorilla/mux"
)

var (
	cosmosPort    = getEnv("COSMOS_PORT", "8888")
	dbHost        = getEnv("INFLUXDB_HOST", "localhost")
	dbPort        = getEnv("INFLUXDB_PORT", "8086")
	dbUsername    = getEnv("INFLUXDB_USERNAME", "root")
	dbPassword    = getEnv("INFLUXDB_PASSWORD", "root")
	dbDatabase    = getEnv("INFLUXDB_DATABASE", "cosmos")
	dbShardConf   = getEnv("INFLUXDB_SHARD_CONF", "./shard_config.json")
	cosmosService = service.NewCosmosService(5)
	newsFeedWorker = worker.NewNewsFeedWorker(cosmosService.LifeTime, 30)
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

// A middleware to serve a cosmos context
//
func serveContext(prev func(
    context.CosmosContext,
    http.ResponseWriter,
    *http.Request)) func(http.ResponseWriter, *http.Request) {
    return (func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        c := context.CosmosContext{
            cosmosService,
            mux.Vars(r) }
        prev(c, w, r)
        return
    })
}

// to create an influxdb client
//
func createInfluxDBClient() (*influxdbc.InfluxDB, error) {
    host := fmt.Sprintf("%s:%s", dbHost, dbPort)
	db := influxdbc.NewInfluxDB(host, dbDatabase, dbUsername, dbPassword)
	file, err := ioutil.ReadFile(dbShardConf)
	if err != nil {
        return db, err
	}

	var conf influxdbc.ShardConfig
	err = json.Unmarshal(file, &conf)
	if err != nil {
        return db, err
	}
    
	_, err = db.CreateDatabase(conf)
	if err != nil {
        return db, err
	} else {
        fmt.Println("[InfluxDB] A database is created.")
    }
    
	return db, nil
}

// to run
//
func run() {
    publicPath := path.Join("telescope", "public")
    
    mux := mux.NewRouter()
    
    mux.HandleFunc("/v1/newsfeeds",
        serveContext(router.GetNewsFeeds)).Methods("GET")
    
    mux.HandleFunc("/v1/planets/{planet}/containers/{container}",
        serveContext(router.GetContainerMetrics)).Methods("GET")
    
    mux.HandleFunc("/v1/planets/{planet}/containers",
        serveContext(router.GetContainersOfPlanet)).Methods("GET")
    
    mux.HandleFunc("/v1/planets/{planet}",
        serveContext(router.GetPlanetMetrics)).Methods("GET")
    
    mux.HandleFunc("/v1/planets",
        serveContext(router.GetPlanets)).Methods("GET")
    
    mux.HandleFunc("/v1/containers",
        serveContext(router.GetContainers)).Methods("GET")
    
    mux.HandleFunc("/v1/planets/{planet}/containers",
        serveContext(router.AddContainersOfPlanet)).Methods("POST")
    
    mux.PathPrefix("/").Handler(http.FileServer(http.Dir(publicPath)))

    middlewares := serveIndexHTML(mux)

    http.Handle("/", middlewares)
    http.ListenAndServe(":" + cosmosPort, nil)
}

func init() {
	db, err := createInfluxDBClient()
    if err != nil {
        fmt.Println(err.Error())
    }
    dao.Initialize(db)
    
	newsFeedWorker.Run()
}

func main() {
	run()
}
