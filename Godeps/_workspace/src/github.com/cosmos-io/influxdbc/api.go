package influxdbc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type InfluxDB struct {
	host     string
	database string
	username string
	password string
}

func escape(q string) string {
	return url.QueryEscape(q)
}

func NewInfluxDB(host, database, username, password string) *InfluxDB {
	return &InfluxDB{host: host, database: database, username: username, password: password}
}

func (db *InfluxDB) SeriesURL(timePrecision string) string {
	return fmt.Sprintf("http://%s/db/%s/series?u=%s&p=%s&time_precision=%s", db.host, escape(db.database), escape(db.username), escape(db.password), escape(timePrecision))
}

func (db *InfluxDB) QueryURL(query, timePrecision string) string {
	return fmt.Sprintf("http://%s/db/%s/series?u=%s&p=%s&q=%s&time_precision=%s", db.host, escape(db.database), escape(db.username), escape(db.password), escape(query), escape(timePrecision))
}

func (db *InfluxDB) WriteSeries(s []*Series, tp string) error {
	url := db.SeriesURL(tp)
	_, err := PostStruct(url, s)
	return err
}

func (db *InfluxDB) Query(query, tp string) ([]*Series, error) {
	url := db.QueryURL(query, tp)
	fmt.Println(url)
	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)
	fmt.Printf("\n%v\n", result.StatusCode)
	fmt.Println(buf.String())

	if result.StatusCode != 200 {
		return nil, errors.New(buf.String())
	}
	var series []*Series
	err = json.Unmarshal(buf.Bytes(), &series)
	if err != nil {
		return nil, err
	}
	return series, nil
}

func PostStruct(url string, reqStruct interface{}) (string, error) {
	marshalled, err := json.Marshal(reqStruct)
	//marshalled = bytes.ToLower(marshalled) commented for case sensitive
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(marshalled)
	result, err := http.Post(url, "application/json", buf)
	if err != nil {
		return "", err
	}
	defer result.Body.Close()
	result_buf := new(bytes.Buffer)
	result_buf.ReadFrom(result.Body)
	if result.StatusCode != 200 {
		return "", errors.New(result_buf.String())
	}
	return result_buf.String(), nil
}
