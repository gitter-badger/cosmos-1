package main

import (
	"encoding/json"
	"fmt"

	"github.com/brannpark/cosmos/model"
	"github.com/cosmos-io/influxdbc"
)

func ConvertToContainerSeries(token, planet string, body []byte) ([]*influxdbc.Series, error) {
	var containers []*model.Container
	err := json.Unmarshal(body, &containers)
	if err != nil {
		return nil, err
	}

	result := make([]*influxdbc.Series, len(containers))
	for i, cont := range containers {

		series := influxdbc.NewSeries(fmt.Sprintf("%s_%s_%s", token, planet, cont.Id))

		cols := []string{"id", "command", "image", "name", "port", "status", "cpu_usage", "mem_usage"}
		points := make([][]interface{}, 1)

		series.Columns = cols
		series.Points = points

		points[0] = make([]interface{}, len(cols))
		points[0][0] = cont.Id
		points[0][1] = cont.Command
		points[0][2] = cont.Image
		points[0][3] = cont.Names[0]
		points[0][4] = cont.Ports[0].Description()
		points[0][5] = cont.Status
		points[0][6] = cont.Stats.Cpu.TotalUtilization
		points[0][7] = cont.Stats.Memory.Usage

		fmt.Println(*series)
		result[i] = series
	}

	return result, nil
}

func ConvertToPlanetSeries(token string, body []byte) (*influxdbc.Series, error) {
	var raw map[string]interface{}
	err := json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}

	cols := make([]string, len(raw))
	points := make([][]interface{}, 1)

	series := influxdbc.NewSeries(token)
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

func ConvertFromSeries(series []*influxdbc.Series) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, s := range series {
		for _, point := range s.Points {
			m := make(map[string]interface{})
			for i, val := range point {
				m[s.Columns[i]] = val
			}
			result = append(result, m)
		}
	}
	return result
}
