package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func getBodyFromRequest(req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func addContainers(r render.Render, params martini.Params, req *http.Request) {
	req.ParseForm()

	body, err := getBodyFromRequest(req)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}
	series, err := ConvertToContainerSeries(params["planet"], body)

	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	s := make([]*influxdbc.Series, 1)
	s[0] = series
	err = db.WriteSeries(s, "")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, series)
}

func getContainers(r render.Render, params martini.Params) {
	series, err := db.Query(fmt.Sprintf("SELECT * FROM containers WHERE planet='%s'", params["planet"]), "")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, series)
}

func getPlanets(r render.Render) {
	series, err := db.Query("SELECT * FROM planets", "")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, series)
}

func addPlanets(r render.Render, req *http.Request) {
	req.ParseForm()
	body, err := getBodyFromRequest(req)
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
	err = db.WriteSeries(s, "")
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, series)
}

// func getContainer(r render.Render, params martini.Params) {
// 	series, err := db.Query(fmt.Sprintf("SELECT * FROM containers WHERE host='%s' AND name='%s'", params["host"], params["name"]), "")
// 	if err != nil {
// 		r.JSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	r.JSON(http.StatusOK, series)
//}
