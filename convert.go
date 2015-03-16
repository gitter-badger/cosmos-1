package main

import (
	"encoding/json"

	"github.com/cosmos-io/influxdbc"
)

func ConvertToContainerSeries(planet string, body []byte) (*influxdbc.Series, error) {
	var raw []map[string]interface{}
	err := json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}

	cols := make([]string, len(raw[0])+1)
	cols[0] = "planet"

	points := make([][]interface{}, len(raw))

	series := influxdbc.NewSeries("containers")
	series.Columns = cols
	series.Points = points

	for i, r := range raw {
		j := 1
		points[i] = make([]interface{}, len(cols))
		points[i][0] = planet
		for k, v := range r {
			if i == 0 {
				cols[j] = k
			}
			points[i][j] = v
			j += 1
		}
	}

	return series, nil
}

func ConvertToPlanetSeries(body []byte) (*influxdbc.Series, error) {
	var raw map[string]interface{}
	err := json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}

	cols := make([]string, len(raw))
	points := make([][]interface{}, 1)

	series := influxdbc.NewSeries("planets")
	series.Columns = cols
	series.Points = points
	points[0] = make([]interface{}, len(raw))

	i := 0
	for k, v := range raw {
		cols[i] = k
		points[0][i] = v
		i = i + 1
	}

	return series, nil
}
