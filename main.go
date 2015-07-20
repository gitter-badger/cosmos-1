package main

import (
	"os"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
    "runtime"
	"io/ioutil"
	"net/http"
	"encoding/json"

	"github.com/cosmoshq/cosmos/context"
	"github.com/cosmoshq/cosmos/influxdb"
	"github.com/cosmoshq/cosmos/route"

	"github.com/gorilla/mux"
)

var (
    port = getEnv("PORT", "8888")

	influxdbHost     = getEnv("INFLUXDB_HOST", "localhost")
	influxdbPort     = getEnv("INFLUXDB_PORT", "8086")
	influxdbUsername = getEnv("INFLUXDB_USERNAME", "root")
	influxdbPassword = getEnv("INFLUXDB_PASSWORD", "root")
	influxdbDatabase = getEnv("INFLUXDB_DATABASE", "cosmos")
	influxdbClient   = newInfluxDB()

    telescopePath = getTelescopePath()
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

// to get an environment variable if it exists or default value
//
func getTelescopePath() string {
    telescopePath := getEnv("TELESCOPE_PATH", "")
    if telescopePath == "" {
        _, current, _, _ := runtime.Caller(1)
        telescopePath = path.Join(path.Dir(current), "telescope")
    }
    return telescopePath
}

// A middleware to distingush text/html content type and others
//
func serveMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        reg := regexp.MustCompile("^.*\\.(jpg|jpeg|png|gif|css|js|html|xml|json|map|txt|ico)$")
		if reg.MatchString(r.URL.Path) {
			fp := path.Join(telescopePath, "public", r.URL.Path)
			http.ServeFile(w, r, fp)
			return
		}

		accept := strings.ToLower(r.Header.Get("Accept"))
		if strings.Contains(accept, "text/html") {
			fp := path.Join(telescopePath, "public", "index.html")
			http.ServeFile(w, r, fp)
			return
		}

        next.ServeHTTP(w, r)
        return
	})
}

// A middleware to serve a cosmos context
//
type MuxHandler func(context.Context, http.ResponseWriter, *http.Request) (int, map[string]interface{})

func serveContext(next MuxHandler) func(http.ResponseWriter, *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		queryParams := r.URL.Query()
		c := context.Context{
			influxdbClient,
			mux.Vars(r),
			body,
			queryParams,
		}

		var js []byte
		var length string

		status, res := next(c, w, r)
		if res == nil {
			js = []byte("")
			length = "0"
		} else {
			js, _ = json.Marshal(res)
			length = strconv.Itoa(len(js))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", length)
		w.WriteHeader(status)
		w.Write(js)

		return
	})
}

// to create an influxdb client
//
func newInfluxDB() *influxdb.InfluxDB {
	rpOneMonth := &influxdb.RetentionPolicy{
		Name:        "30d.metrics",
		Duration:    "30d",
		Replication: 1,
		Default:     true,
	}

	retentionPolicies := make([]*influxdb.RetentionPolicy, 1)
	retentionPolicies[0] = rpOneMonth

	config := influxdb.Config{
		Host:              influxdbHost,
		Port:              influxdbPort,
		Username:          influxdbUsername,
		Password:          influxdbPassword,
		Database:          influxdbDatabase,
		RetentionPolicies: retentionPolicies,
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
	mux := mux.NewRouter()

	mux.HandleFunc("/metrics",
		serveContext(route.PostMetrics)).Methods("POST")

	mux.HandleFunc("/planets",
		serveContext(route.GetPlanets)).Methods("GET")

	mux.HandleFunc("/containers",
		serveContext(route.GetContainers)).Methods("GET")

	mux.HandleFunc("/metrics",
		serveContext(route.GetMetrics)).Methods("GET")

	http.Handle("/", serveMiddleware(mux))
	http.ListenAndServe(":" + port, nil)
}

func init() {
}

func main() {
	run()
}
