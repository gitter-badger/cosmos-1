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

func getAllPlanets(r render.Render, req *http.Request) {
	token := getToken(req)

	series, err := logDb.Query(fmt.Sprintf("SELECT * FROM /^%s\\.[^\\.]+$/ LIMIT 1", token), "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	planets := make([]map[string]interface{}, len(series))
	for i, s := range series {
		planets[i] = make(map[string]interface{})
		planets[i]["Name"] = strings.Split(s.Name, ".")[1]
	}

	r.JSON(http.StatusOK, planets)
}

func getAllContainers(r render.Render, req *http.Request) {
	token := getToken(req)

	series, err := logDb.Query(fmt.Sprintf("SELECT txt_value, num_value FROM /^min\\.%s\\..*/ LIMIT 1", token), "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, ConvertFromContainerSeries("", series))
}

func addContainersOfPlanet(r render.Render, params martini.Params, req *http.Request) {
	req.ParseForm()
	token := getToken(req)
	planet := params["planetName"]

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
	planetSeries := influxdbc.NewSeries(fmt.Sprintf("%s.%s", token, planet), "value")
	planetSeries.AddPoint("")
	series = append(series, planetSeries)

	err = logDb.WriteSeries(series, "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.Status(http.StatusOK)
}

func getContainersOfPlanet(r render.Render, params martini.Params, req *http.Request) {
	token := getToken(req)
	planet := params["planetName"]

	dbQuery := fmt.Sprintf("SELECT txt_value, num_value FROM /^min\\.%s\\.%s\\./ LIMIT 1", token, planet)
	series, err := logDb.Query(dbQuery, "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, ConvertFromContainerSeries(planet, series))
}

func getContainerInfo(r render.Render, params martini.Params, req *http.Request) {
	token := getToken(req)

	planet := params["planetName"]
	containerName := strings.Replace(params["containerName"], ".", "_", -1)

	seriesName := MakeContainerSeriesName(token, planet, containerName)

	dbQuery := fmt.Sprintf("SELECT txt_value, num_value FROM /^min\\.%s\\./ LIMIT 10", seriesName)
	series, err := logDb.Query(dbQuery, "s")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, ConvertFromContainerInfoSeries(containerName, series))

}
