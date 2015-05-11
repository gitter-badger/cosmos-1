package router

import (
    "fmt"
    "net/http"
    "encoding/json"
    
    "github.com/cosmos-io/cosmos/context"
    "github.com/cosmos-io/cosmos/model"
)

var (
    types = map[string]bool {
        "cpu": true,
        "memory": true,
    }
)

func PostMetrics(
    c context.Context,
    w http.ResponseWriter,
    r *http.Request) (int, map[string]interface{}) {
    var metrics *model.MetricsParam
    err := json.Unmarshal(c.Body, &metrics)

    var status int
    var res map[string]interface{}
    
    if err != nil {
        status = http.StatusBadRequest
        res = map[string]interface{} { "error": err.Error() }
        return status, res
    }

    c.InfluxDB.WriteMetrics(metrics)

    status = http.StatusOK
    
    return status, res
}

func GetMetrics(
    c context.Context,
    w http.ResponseWriter,
    r *http.Request) (int, map[string]interface{}) {
    t := c.GetQueryParam("type", "")
    planet := c.GetQueryParam("planet", "")
    container := c.GetQueryParam("container", "")

    var status int
    var res map[string]interface{}

    if types[t] == false {
        status = http.StatusBadRequest
        err := fmt.Sprintf("%s type is not supported.", t)
        res = map[string]interface{} { "error": err }
        return status, res
    }

    if planet == "" {
        status = http.StatusBadRequest
        err := fmt.Sprintf("planet is empty.")
        res = map[string]interface{} { "error": err }
        return status, res
    }

    var metrics interface{}
    var err error

    if container == "" {
        metrics, err = c.InfluxDB.QueryPlanetMetrics(planet, t)
    } else {
        metrics, err = c.InfluxDB.QueryContainerMetrics(planet, container, t)
    }

    if err != nil {
        status = http.StatusInternalServerError
        res := map[string]interface{} { "error": err.Error() }
        return status, res
    }

    status = http.StatusOK
    res = map[string]interface{} {
        "data": metrics,
    }

    return status, res
}
