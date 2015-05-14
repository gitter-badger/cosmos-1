package influxdb

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/cosmos-io/cosmos/model"

	"github.com/influxdb/influxdb/client"
)

type Config struct {
	Host              string
	Port              string
	Username          string
	Password          string
	Database          string
	RetentionPolicies []*RetentionPolicy
}

type RetentionPolicy struct {
	Name        string
	Duration    string
	Replication int
	Default     bool
}

type InfluxDB struct {
	client   *client.Client
	database string
}

var (
	cosmos = "cosmos"
)

func New(config Config) (*InfluxDB, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s:%s", config.Host, config.Port))
	if err != nil {
		return nil, err
	}

	conf := client.Config{
		URL:      *u,
		Username: config.Username,
		Password: config.Password,
	}

	client, err := client.NewClient(conf)
	if err != nil {
		return nil, err
	}

	db := &InfluxDB{
		client:   client,
		database: config.Database,
	}

	// Creating a database
	err = db.createDatabase()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("database has been created")
	}
	// Creating retention policies
	err = db.createRetentionPolicy(config.RetentionPolicies)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("retention policies has been created")
	}

	return db, nil
}

// queryDB convenience function to query the database
func (db *InfluxDB) Query(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db.database,
	}
	if response, err := db.client.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
}

func (db *InfluxDB) createDatabase() error {
	_, err := db.Query(fmt.Sprintf("create database %s", db.database))
	return err
}

func (db *InfluxDB) createRetentionPolicy(retentionPolicies []*RetentionPolicy) error {
	if retentionPolicies == nil || len(retentionPolicies) == 0 {
		return nil
	}

	for _, rp := range retentionPolicies {
		cmd := fmt.Sprintf("CREATE RETENTION POLICY \"%s\" ON %s DURATION %s REPLICATION %d", rp.Name, db.database, rp.Duration, rp.Replication)
		if rp.Default == true {
			cmd += " DEFAULT"
		}
		_, err := db.Query(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *InfluxDB) WriteMetrics(metrics *model.MetricsParam) {
	planet := metrics.Planet
	if planet == "" {
		return
	}

	index := 0
	sampleSize := len(metrics.Containers)
	pts := make([]client.Point, sampleSize*2) // cpu, memory
	for i := 0; i < sampleSize; i++ {
		container := metrics.Containers[i]
		if container.Container == "" {
			continue
		}

		// cpu
		pts[index] = client.Point{
			Name: "cpu",
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
			Name: "memory",
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

func (db *InfluxDB) QueryPlanetMetrics(planet string, t string) (interface{}, error) {
	cmd := fmt.Sprintf("SELECT max(value) FROM %s WHERE cosmos = '%s' AND planet = '%s' AND time > now() - %s group by container, time(%s), value fill(0)",
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
