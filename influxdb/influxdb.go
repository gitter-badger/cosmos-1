package influxdb

import (
    "fmt"
    "net/url"
    
    "github.com/influxdb/influxdb/client"
)

type InfluxDB struct {
    client *client.Client
}

func NewClient(host string,
    port string,
    username string,
    password string) (*InfluxDB, error) {
    u, err := url.Parse(fmt.Sprintf("http://%s:%s", host, port))
    if err != nil {
        return nil, err
    }

    conf := client.Config {
        URL:      *u,
        Username: username,
        Password: password,
    }

    con, err := client.NewClient(conf)
    if err != nil {
        return nil, err
    }
    
    influxdb := InfluxDB {
        client: con,
    }
    return &influxdb, nil
}

/*u, err := url.Parse(fmt.Sprintf("http://%s:%s", "localhost", "8086"))
    if err != nil {
        fmt.Println(err)
    }

    conf := client.Config{
        URL:      *u,
        Username: "root",
        Password: "root",
    }

    con, err := client.NewClient(conf)
    if err != nil {
        fmt.Println(err)
    }

    sampleSize := 1
    pts := make([]client.Point, sampleSize)

    pts[0] = client.Point {
        Name: "cpu",
        Tags: map[string]string {
            "region": "useast",
            "host": "server02",
        },
        Fields: map[string]interface{}{
            "value": 100,
        },
        Timestamp: time.Now(),
        Precision: "s",
    }

    bps := client.BatchPoints{
        Points:          pts,
        Database:        "cosmos",
        RetentionPolicy: "default",
    }

    _, err = con.Write(bps)
    if err != nil {
        fmt.Println(err)
    }*/
