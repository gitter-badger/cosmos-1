package main

import (
	"fmt"
	"net/http"

	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func addContainers(r render.Render, params martini.Params, req *http.Request) {
	req.ParseForm()
	token := req.URL.Query().Get("token")
	if token == "" {
		token = "default"
	}
	planet := params["planet"]

	body, err := GetBodyFromRequest(req)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}
	series, err := ConvertToContainerSeries(token, planet, body)

	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	s := make([]*influxdbc.Series, 1)
	s[0] = series
	err = db.WriteSeries(s, "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, series)
}

func getContainers(r render.Render, params martini.Params, req *http.Request) {
	stime := req.URL.Query().Get("stime")
	etime := req.URL.Query().Get("etime")
	token := req.URL.Query().Get("token")
	if token == "" {
		token = "default"
	}
	planet := params["planet"]
	containerSeriesName := GenerateContainerSeriesName(token, planet)

	dbQuery := fmt.Sprintf("SELECT * FROM %s WHERE time >= %ss and time <= %ss", containerSeriesName, stime, etime), "s"
	series, err := db.Query(dbQuery)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, series)
}

func getPlanets(r render.Render) {
	series, err := db.Query("SELECT * FROM planets", "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, series)
}

func addPlanets(r render.Render, req *http.Request) {
	req.ParseForm()
	body, err := GetBodyFromRequest(req)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}
	series, err := ConvertToPlanetSeries(body)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	fmt.Println(series)
	s := make([]*influxdbc.Series, 1)
	s[0] = series
	err = db.WriteSeries(s, "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, series)
}
