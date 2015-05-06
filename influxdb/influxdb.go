package influxdb

import (
    "fmt"
    "log"
    "time"
    "net/url"

    "github.com/cosmos-io/cosmos/model"
    
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
    cosmosName = "cosmos"
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

func (i *InfluxDB) WriteMetrics(metrics *model.MetricsParam) {
    cosmos := cosmosName
    planet := metrics.Planet
    if planet == "" { return }

    index := 0
    sampleSize := len(metrics.Containers)
    pts := make([]client.Point, sampleSize)
    for i := 0; i < sampleSize; i++ {
        container := metrics.Containers[i]
        if container.Container == "" { continue }
        pts[index] = client.Point {
            Name: "cpu",
            Tags: map[string]string {
                "cosmos": cosmos,
                "planet": planet,
                "container": container.Container,
            },
            Fields: map[string]interface{} {
                "value": container.Cpu,
            },
            Timestamp: time.Now(),
            Precision: "s",
        }
        index += 1
    }

    bps := client.BatchPoints {
        Points:          pts,
        Database:        databaseName,
        RetentionPolicy: retentionPolicy,
    }

    _, err := i.client.Write(bps)
    if err != nil {
        log.Println(err)
    }
}

func (i *InfluxDB) QueryPlanets() ([]string, error) {
    c := fmt.Sprintf("SHOW TAG VALUES WITH KEY = %s WHERE cosmos = '%s'",
        "planet",
        cosmosName,
    )

    result, err := queryDB(i.client, c)
    if err != nil {
        return nil, err
    }

    if len(result) == 0 || len(result[0].Series) == 0 {
        // no value
    }

    values := result[0].Series[0].Values
    length := len(values)

    planets := make([]string, length)
    for i := 0; i < length; i++ {
        planets[i] = values[i][0].(string)
    }
    return planets, nil
}

func (i *InfluxDB) QueryContainers(planet string) ([]string, error) {
    c := fmt.Sprintf("SHOW TAG VALUES WITH KEY = %s WHERE cosmos = '%s' AND planet = '%s'",
        "container",
        cosmosName,
        planet,
    )

    result, err := queryDB(i.client, c)
    if err != nil {
        return nil, err
    }

    if len(result) == 0 || len(result[0].Series) == 0 {
        // no value
    }

    values := result[0].Series[0].Values
    length := len(values)

    containers := make([]string, length)
    for i := 0; i < length; i++ {
        containers[i] = values[i][0].(string)
    }
    return containers, nil
}
