package main

import (
	"fmt"
	"os"

	"github.com/attilaolah/strict"
	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	db *influxdbc.InfluxDB

	cosmosPort = getEnv("COSMOS_PORT", "8080")
	dbHost     = getEnv("INFLUXDB_HOST", "localhost")
	dbPort     = getEnv("INFLUXDB_PORT", "8086")
	dbUsername = getEnv("INFLUXDB_USERNAME", "root")
	dbPassword = getEnv("INFLUXDB_PASSWORD", "root")
	dbName     = getEnv("INFLUXDB_DATABASE", "cosmos")
)

func getEnv(key, defVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defVal
	} else {
		return val
	}
}

func startServer() {
	db = influxdbc.NewInfluxDB(fmt.Sprintf("%s:%s", dbHost, dbPort), dbName, dbUsername, dbPassword)
	m := martini.Classic()

	m.Handlers(
		martini.Logger(),
		martini.Static("public"),
		strict.Strict,
		render.Renderer(),
	)

	m.Group("/v1", func(r martini.Router) {
		r.Post("/:host/containers", strict.Accept("application/json"), strict.ContentType("application/json"), postContainers)
		r.Get("/:host/containers", strict.Accept("application/json"), getContainers)
		r.Get("/:host/containers/:name", strict.Accept("application/json"), getContainer)
	})

	if cosmosPort == "" {
		m.RunOnAddr(":3000")
	} else {
		m.RunOnAddr(":" + cosmosPort)
	}
}

func main() {
	startServer()
}
