package main

import (
	"fmt"
	"net/http"
	"strings"

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

	series, err := ConvertToContainerSeries(token, planet, body)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	err = logDb.WriteSeries(series, "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	// cSeries := make([]*influxdbc.Series, 1)
	// names := make([]string, len(series))

	// cs := influxdbc.NewSeries(fmt.Sprintf("%s_%s", token, planet), "containers")
	// cSeries[0] = cs
	// for i, s := range series {
	// 	names[i] = strings.Split(s.Name, "_")[2]
	// }
	// cs.AddPoint(strings.Join(names, ","))

	// err = logDb.WriteSeries(cSeries, "s")
	// if err != nil {
	// 	fmt.Println(err)
	// 	r.JSON(http.StatusInternalServerError, err)
	// 	return
	// }

	r.Status(http.StatusOK)
}

func getContainers(r render.Render, params martini.Params, req *http.Request) {
	interval := req.URL.Query().Get("interval")
	token := getToken(req)

	planet := params["planet"]
	//seriesName := GenerateContainerSeriesName(token, planet)

	dbQuery := fmt.Sprintf("SELECT * FROM /%s_%s_[^_]+$/ WHERE time > now() - %s LIMIT 1", token, planet, interval)
	series, err := logDb.Query(dbQuery, "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	containers := make([]map[string]interface{}, len(series))
	for i, s := range series {
		containers[i] = make(map[string]interface{})
		containers[i]["name"] = strings.Split(s.Name, "_")[2]
	}

	r.JSON(http.StatusOK, containers)
}

func getContainerInfo(r render.Render, params martini.Params, req *http.Request) {
	interval := req.URL.Query().Get("interval")
	token := getToken(req)

	planet := params["planet"]
	container := params["container"]
	seriesName := GenerateContainerInfoSeriesName(token, planet, container)

	dbQuery := fmt.Sprintf("SELECT * FROM %s WHERE time > now() - %s", seriesName, interval)
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

// func addPlanets(r render.Render, req *http.Request) {
// 	req.ParseForm()
// 	token := getToken(req)

// 	body, err := GetBodyFromRequest(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		r.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	series, err := ConvertToPlanetSeries(token, body)
// 	if err != nil {
// 		fmt.Println(err)
// 		r.JSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	fmt.Println(series)
// 	s := make([]*influxdbc.Series, 1)
// 	s[0] = series
// 	err = logDb.WriteSeries(s, "s")
// 	if err != nil {
// 		fmt.Println(err)
// 		r.JSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	r.Status(http.StatusOK)
// }
