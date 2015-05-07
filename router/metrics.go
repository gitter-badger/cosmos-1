package router

import (
    "fmt"
    "strconv"
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
    r *http.Request) {
    var metrics *model.MetricsParam
    err := json.Unmarshal(c.Body, &metrics)
    if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        contentLength := strconv.Itoa(len(js))
        w.WriteHeader(http.StatusInternalServerError)
        w.Header().Set("Content-Length", contentLength)
        w.Write(js)
		return
    }

    c.InfluxDB.WriteMetrics(metrics)

    content := []byte("")
    contentLength := strconv.Itoa(len(content))
    w.Header().Set("Content-Length", contentLength)
    w.Write(content)
}

func GetMetrics(
    c context.Context,
    w http.ResponseWriter,
    r *http.Request) {
    t := c.GetQueryParam("type", "")
    planet := c.GetQueryParam("planet", "")
    container := c.GetQueryParam("container", "")

    if types[t] == false {
        err := fmt.Sprintf("%s type is not supported.", t)
        res := map[string]string { "error": err }
        js, _ := json.Marshal(res)
        contentLength := strconv.Itoa(len(js))
        w.WriteHeader(http.StatusBadRequest)
        w.Header().Set("Content-Length", contentLength)
        w.Write(js)
        return
    }

    metrics, err := c.InfluxDB.QueryMetrics(planet, container, t)
    if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        contentLength := strconv.Itoa(len(js))
        w.WriteHeader(http.StatusBadRequest)
        w.Header().Set("Content-Length", contentLength)
        w.Write(js)
        return
    }

    res := map[string]interface{} {
        "data": metrics,
    }

    js, _ := json.Marshal(res)
    contentLength := strconv.Itoa(len(js))
    w.Header().Set("Content-Length", contentLength)
    w.Write(js)
}
