package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func getToken(req *http.Request) string {
	token := req.URL.Query().Get("token")
	if token == "" {
		token = "default"
	}
	return token
}

func addContainers(r render.Render, params martini.Params, req *http.Request) {
	req.ParseForm()
	token := getToken(req)
	planet := params["planet"]

	body, err := GetBodyFromRequest(req)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	// Container metrics
	series, err := ConvertToContainerSeries(token, planet, body)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	// Host metrics
	planetSeries := influxdbc.NewSeries(fmt.Sprintf("%s_%s", token, planet), "desc")
	planetSeries.AddPoint("test")
	series = append(series, planetSeries)

	err = logDb.WriteSeries(series, "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.Status(http.StatusOK)
}

func getContainers(r render.Render, params martini.Params, req *http.Request) {
	ttl := req.URL.Query().Get("ttl")
	token := getToken(req)

	planet := params["planet"]

	dbQuery := fmt.Sprintf("SELECT * FROM /%s_%s_[^_]+$/ WHERE time > now() - %s LIMIT 1", token, planet, ttl)
	series, err := logDb.Query(dbQuery, "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, ConvertFromSeries(series))
}

func getContainerInfo(r render.Render, params martini.Params, req *http.Request) {
	interval := req.URL.Query().Get("interval")
	token := getToken(req)

	planet := params["planet"]
	container := params["container"]
	seriesName := fmt.Sprintf("%s_%s_%s", token, planet, container)

	dbQuery := fmt.Sprintf("SELECT mean(cpu_usage) as cpu_usage, mean(mem_usage) as mem_usage FROM %s GROUP BY time(%s) LIMIT 10", seriesName, interval)
	series, err := logDb.Query(dbQuery, "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, ConvertFromSeries(series))

}

func getPlanets(r render.Render, req *http.Request) {
	token := getToken(req)

	series, err := logDb.Query(fmt.Sprintf("SELECT * FROM /%s_[^_]+$/ LIMIT 1", token), "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	planets := make([]map[string]interface{}, len(series))
	for i, s := range series {
		planets[i] = make(map[string]interface{})
		planets[i]["name"] = strings.Split(s.Name, "_")[1]
	}

	r.JSON(http.StatusOK, planets)
}
