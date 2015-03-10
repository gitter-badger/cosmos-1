package main

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos-io/influxdbc"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

func postContainers(r render.Render, params martini.Params, req *http.Request) {
	req.ParseForm()

	var raw []map[string]interface{}
	err := json.Unmarshal([]byte(req.FormValue("data")), &raw)
	if err != nil {
		fmt.Printf("%s\n", err)
		r.JSON(500, err)
		return
	}

	cols := make([]string, len(raw[0])+1)
	cols[0] = "host"

	points := make([][]interface{}, len(raw))

	series := influxdbc.NewSeries("containers")
	series.Columns = cols
	series.Points = points

	for i, r := range raw {
		j := 1
		points[i] = make([]interface{}, len(cols))
		points[i][0] = params["host"]
		for k, v := range r {
			if i == 0 {
				cols[j] = k
			}
			points[i][j] = v
			j += 1
		}
	}

	s := make([]influxdbc.Series, 1)
	s[0] = *series
	err = db.WriteSeries(s, "")
	if err != nil {
		r.JSON(500, err)
		return
	}

	//res := NewResJson()
	//res["host"] = hostname
	r.JSON(200, series)
}

func getContainers(r render.Render, params martini.Params, req *http.Request) {
	req.ParseForm()

	results, err := db.Query("SELECT * FROM containers", "")
	if err != nil {
		fmt.Println("Error")
		fmt.Printf("%s\n", err)
		r.JSON(500, err)
		return
	}

	r.JSON(200, results)
}
