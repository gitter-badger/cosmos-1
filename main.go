package main

import (
    "os"
    "log"
    "path"
    "strings"
    "strconv"
    "io/ioutil"
    "net/http"
    "encoding/json"

    "github.com/cosmos-io/cosmos/context"
    "github.com/cosmos-io/cosmos/influxdb"
    "github.com/cosmos-io/cosmos/router"
    
    "github.com/gorilla/mux"
)

var (
    cosmosPort       = getEnv("COSMOS_PORT", "8888")
    
	influxdbHost     = getEnv("INFLUXDB_HOST", "localhost")
	influxdbPort     = getEnv("INFLUXDB_PORT", "8086")
	influxdbUsername = getEnv("INFLUXDB_USERNAME", "root")
	influxdbPassword = getEnv("INFLUXDB_PASSWORD", "root")
	influxdbDatabase = getEnv("INFLUXDB_DATABASE", "cosmos")
    influxdbClient   = newInfluxDB()
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
func serveContext(
    next func(
        context.Context,
        http.ResponseWriter,
        *http.Request)(int, map[string]interface{})) func(http.ResponseWriter, *http.Request) {
    return (func(w http.ResponseWriter, r *http.Request) {
        body, _ := ioutil.ReadAll(r.Body)
        queryParams := r.URL.Query()
        c := context.Context {
            influxdbClient,
            mux.Vars(r),
            body,
            queryParams,
        }
        
        status, res := next(c, w, r)
        js, _ := json.Marshal(res)
        len := strconv.Itoa(len(js))

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Content-Length", len)
        w.WriteHeader(status)
        w.Write(js)
        
        return
    })
}

// to create an influxdb client
//
func newInfluxDB() *influxdb.InfluxDB {
    config := influxdb.Config {
        Host: influxdbHost,
        Port: influxdbPort,
        Username: influxdbUsername,
        Password: influxdbPassword,
        Database: influxdbDatabase,
    }
    
    db, err := influxdb.New(config)
    if err != nil {
        log.Println(err)
    }
    
    return db
}

// to run
//
func run() {
    publicPath := path.Join("telescope", "public")

    mux := mux.NewRouter()

    mux.HandleFunc("/metrics",
        serveContext(router.PostMetrics)).Methods("POST")

    mux.HandleFunc("/planets",
        serveContext(router.GetPlanets)).Methods("GET")

    mux.HandleFunc("/containers",
        serveContext(router.GetContainers)).Methods("GET")

    mux.HandleFunc("/metrics",
        serveContext(router.GetMetrics)).Methods("GET")

    /*mux.HandleFunc("/v1/newsfeeds",
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
        serveContext(router.AddContainersOfPlanet)).Methods("POST")*/

    mux.PathPrefix("/").Handler(http.FileServer(http.Dir(publicPath)))

    middlewares := serveIndexHTML(mux)

    http.Handle("/", middlewares)
    http.ListenAndServe(":" + cosmosPort, nil)
}

func init() {
    
}

func main() {
	run()
}
