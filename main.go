package main

import (
	"fmt"
	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"os"
)

var (
	db   *influxdbc.InfluxDB
	host = os.Getenv("INFLUXDB_HOST")
	port = os.Getenv("INFLUXDB_PORT")
)

func startServer() {
	db = influxdbc.NewInfluxDB(fmt.Sprintf("%s:%d", host, port), "cosmos", "root", "root")
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

	m.Run()
}

func main() {
	startServer()
}
