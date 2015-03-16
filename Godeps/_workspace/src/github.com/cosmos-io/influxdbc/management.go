package influxdbc

import (
	"fmt"
	"net/http"
)

type Space struct {
	Name              string
	RetentionPolicy   string
	ShardDuration     string
	Regex             string
	ReplicationFactor int
	Split             int
}

type ShardSpace map[string][]*Space

func (db *InfluxDB) CreateDatabase(shard ShardSpace) (string, error) {
	url := fmt.Sprintf("http://%s/cluster/database_configs/%s?u=%s&p=%s", db.host, db.database, db.username, db.password)
	response, err := PostStruct(url, shard)
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
