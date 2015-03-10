package main

import (
	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	db *influxdbc.InfluxDB
)

func startServer() {
	db = influxdbc.NewInfluxDB("localhost:8086", "cosmos", "root", "root")
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
