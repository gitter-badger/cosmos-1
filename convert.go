package main

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos-io/influxdbc"
)

func GenerateContainerInfoSeriesName(token, planet, container string) string {
	return fmt.Sprintf("%s_%s_%s", token, planet, container)
}

func GenerateContainerSeriesName(token, planet string) string {
	return fmt.Sprintf("%s_%s", token, planet)
}

func ConvertToContainerSeries(token, planet string, body []byte) ([]*influxdbc.Series, error) {
	var containers map[string]map[string]interface{}
	err := json.Unmarshal(body, &containers)
	if err != nil {
		return nil, err
	}

	result := make([]*influxdbc.Series, 0)
	for name, infos := range containers {
		series := influxdbc.NewSeries(GenerateContainerInfoSeriesName(token, planet, name))

		cols := make([]string, len(infos))
		points := make([][]interface{}, 1)
		points[0] = make([]interface{}, len(cols))

		series.Columns = cols
		series.Points = points

		i := 0
		for k, v := range infos {
			cols[i] = k
			points[0][i] = v
			i += 1
		}
		fmt.Println(*series)
		result = append(result, series)
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
