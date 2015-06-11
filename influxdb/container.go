package influxdb

import (
	"fmt"
	"log"
	"time"

	"github.com/cosmos-io/cosmos/model"
	"github.com/influxdb/influxdb/client"
)

func (db *InfluxDB) QueryContainers(planet string) ([]string, error) {
	cmd := fmt.Sprintf("SHOW TAG VALUES WITH KEY = %s WHERE cosmos = '%s' AND planet = '%s'",
		"container",
		cosmos,
		planet,
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

	containers := make([]string, length)
	for i := 0; i < length; i++ {
		containers[i] = values[i][0].(string)
	}
	return containers, nil
}

func (db *InfluxDB) QueryContainerMetrics(planet string, container string, t string) (interface{}, error) {
	cmd := fmt.Sprintf("SELECT max(value) FROM %s WHERE cosmos = '%s' AND planet = '%s' AND container = '%s' AND time > now() - %s group by time(%s), value fill(0)",
		t,
		cosmos,
		planet,
		container,
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

	values := result[0].Series[0].Values
	return values, nil
}

func (db *InfluxDB) WriteMetrics(metrics *model.MetricsParam) {
	planet := metrics.Planet
	if planet == "" {
		return
	}

	index := 0
	sampleSize := len(metrics.Containers)
	pts := make([]client.Point, sampleSize * 2) // cpu, memory
	for i := 0; i < sampleSize; i++ {
		container := metrics.Containers[i]
		if container.Container == "" {
			continue
		}

		// cpu
		pts[index] = client.Point{
			Measurement: "cpu",
			Tags: map[string]string{
				"cosmos":    cosmos,
				"planet":    planet,
				"container": container.Container,
			},
			Fields: map[string]interface{}{
				"value": container.Cpu,
			},
			Time:      time.Now(),
			Precision: "s",
		}

		// memory
		pts[index+1] = client.Point{
			Measurement: "memory",
			Tags: map[string]string{
				"cosmos":    cosmos,
				"planet":    planet,
				"container": container.Container,
			},
			Fields: map[string]interface{}{
				"value": container.Memory,
			},
			Time:      time.Now(),
			Precision: "s",
		}

		index += 2 // cpu, memory
	}

	bps := client.BatchPoints{
		Points:   pts,
		Database: db.database,
	}

	_, err := db.client.Write(bps)
	if err != nil {
		log.Println(err)
	}
}

func (db *InfluxDB) QueryFirstContainerMetrics() ([]map[string]string, error) {
	//	time := lastCheckTime.UTC().Format(time.RFC3339Nano)
	cmd := "SELECT * FROM cpu GROUP BY cosmos, planet, container LIMIT 1"

	result, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}

	newContainers := make([]map[string]string, len(result[0].Series))
	for i, s := range result[0].Series {
		newContainers[i] = s.Tags
		newContainers[i]["time"] = s.Values[0][0].(string)
	}

	return newContainers, nil
}

func (db *InfluxDB) QueryContainersInRange(start, end, timeGroup string) (map[string]map[string]string, error) {
	cmd := fmt.Sprintf("SELECT MAX(value) FROM cpu WHERE time > NOW() - %s and time < now() - %s GROUP BY time(%s), cosmos, planet, container LIMIT 1", start, end, timeGroup)
	result, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 || len(result[0].Series) == 0 {
		return nil, nil
	}

	containers := make(map[string]map[string]string)
	for _, s := range result[0].Series {
		key := fmt.Sprintf("%s.%s.%s", s.Tags["cosmos"], s.Tags["planet"], s.Tags["container"])
		containers[key] = s.Tags
	}

	return containers, nil
}
