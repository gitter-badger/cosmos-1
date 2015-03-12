package main

import (
	"encoding/json"

	"github.com/cosmos-io/influxdbc"
)

func ConvertReqBodyToSeries(host string, body []byte) (*influxdbc.Series, error) {
	var raw []map[string]interface{}
	err := json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
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
		points[i][0] = host
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
