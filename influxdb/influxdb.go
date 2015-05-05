package influxdb

import (
    "fmt"
    "log"
    "time"
    "net/url"
    
    "github.com/influxdb/influxdb/client"
)

type Config struct {
    Host string
    Port string
    Username string
    Password string
}

type InfluxDB struct {
    client *client.Client
}

var (
    databaseName = "cosmos"
    retentionPolicy = "default"
)

func New(config Config) (*InfluxDB, error) {
    u, err := url.Parse(fmt.Sprintf("http://%s:%s", config.Host, config.Port))
    if err != nil {
        return nil, err
    }

    conf := client.Config {
        URL:      *u,
        Username: config.Username,
        Password: config.Password,
    }

    con, err := client.NewClient(conf)
    if err != nil {
        return nil, err
    }

    // creating a database
    _, err = queryDB(con, fmt.Sprintf("create database %s", databaseName))
    if err != nil {
        log.Println(err)
    }
    
    influxdb := InfluxDB {
        client: con,
    }
    return &influxdb, nil
}

// queryDB convenience function to query the database
func queryDB(con *client.Client, cmd string) (res []client.Result, err error) {
    q := client.Query {
        Command:  cmd,
        Database: databaseName,
    }
    if response, err := con.Query(q); err == nil {
        if response.Error() != nil {
            return res, response.Error()
        }
        res = response.Results
    }
    return
}

func (i *InfluxDB) WriteExample() {
    sampleSize := 1
    pts := make([]client.Point, sampleSize)

    pts[0] = client.Point {
        Name: "cpu",
        Tags: map[string]string {
            "region": "useast",
            "host": "server01",
        },
        Fields: map[string]interface{} {
            "value": 100,
        },
        Timestamp: time.Now(),
        Precision: "s",
    }

    bps := client.BatchPoints {
        Points:          pts,
        Database:        databaseName,
        RetentionPolicy: retentionPolicy,
    }

    _, err := i.client.Write(bps)
    if err != nil {
        fmt.Println(err) // for test
    }
}
