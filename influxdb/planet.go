package influxdb

import (
    "fmt"
)

func (db *InfluxDB) QueryPlanets() ([]string, error) {
	cmd := fmt.Sprintf("SHOW TAG VALUES WITH KEY = %s WHERE cosmos = '%s'",
		"planet",
		cosmos,
	)

	result, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 || len(result[0].Series) == 0 {
		return nil, nil
	}

	values := result[0].Series[0].Values
	length := len(values)

	planets := make([]string, length)
	for i := 0; i < length; i++ {
		planets[i] = values[i][0].(string)
	}
	return planets, nil
}

func (db *InfluxDB) QueryPlanetMetrics(planet string, t string) (interface{}, error) {
	cmd := fmt.Sprintf("SELECT max(value) FROM %s WHERE cosmos = '%s' AND planet = '%s' AND time > now() - %s group by container, time(%s) fill(0)",
		t,
		cosmos,
		planet,
		"24h",
		"10m",
	)

	result, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 || len(result[0].Series) == 0 {
		return nil, nil
	}

	var values = make(map[string]interface{})
	for i := 0; i < len(result[0].Series); i++ {
		series := result[0].Series[i]
		container := series.Tags["container"]
		values[container] = series.Values
	}

	return values, nil
}
