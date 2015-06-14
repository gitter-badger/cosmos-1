package influxdb

import (
	"fmt"
	"log"
	"net/url"

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
	}

	// Creating retention policies
	err = db.createRetentionPolicy(config.RetentionPolicies)
	if err != nil {
		log.Println(err)
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
