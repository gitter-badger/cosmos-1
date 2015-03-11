package main

import (
	"fmt"
	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"os"
)

var (
	db *influxdbc.InfluxDB

	cosmosPort = os.Getenv("COSMOS_PORT")
	dbHost     = os.Getenv("INFLUXDB_HOST")
	dbPort     = os.Getenv("INFLUXDB_PORT")
	dbUsername = os.Getenv("INFLUXDB_USERNAME")
	dbPassword = os.Getenv("INFLUXDB_PASSWORD")
	dbName     = os.Getenv("INFLUXDB_DATABASE")
)

func startServer() {
	db = influxdbc.NewInfluxDB(fmt.Sprintf("%s:%s", dbHost, dbPort), dbName, dbUsername, dbPassword)
	m := martini.Classic()

	m.Handlers(
		martini.Logger(),
		martini.Static("public"),
	)

	m.Use(render.Renderer())

	m.Group("/v1", func(r martini.Router) {
		r.Post("/:host/containers", postContainers)
		r.Get("/:host/containers", getContainers)
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
