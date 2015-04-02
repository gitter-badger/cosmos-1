package influxdbc

import (
	"fmt"
	"net/http"
)

type Space struct {
	Name              string `json:"name"`
	RetentionPolicy   string `json:"retentionPolicy"`
	ShardDuration     string `json:"shardDuration"`
	Regex             string `json:"regex"`
	ReplicationFactor int    `json:"replicationFactor"`
	Split             int    `json:"split"`
}

type ShardConfig struct {
	Spaces            []*Space `json:"spaces"`
	ContinuousQueries []string `json:"continuousQueries"`
}

func (db *InfluxDB) CreateDatabase(conf ShardConfig) (string, error) {
	url := fmt.Sprintf("http://%s/cluster/database_configs/%s?u=%s&p=%s", db.host, db.database, db.username, db.password)
	response, err := PostStruct(url, conf)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (db *InfluxDB) DeleteDatabase(database string) error {
	url := fmt.Sprintf("http://%s/db/%s?u=%s&p=%s", db.host, db.database, db.username, db.password)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	result, _ := http.DefaultClient.Do(req)
	defer result.Body.Close()
	return nil
}
