package router

import (
    "log"
    "net/http"
    "encoding/json"
    
    "github.com/cosmos-io/cosmos/context"
    "github.com/cosmos-io/cosmos/model"
)

func PostMetrics(
    c context.CosmosContext,
    w http.ResponseWriter,
    r *http.Request) {
    var metrics *model.Metrics
    err := json.Unmarshal(c.Body, &metrics)
    if err != nil {
        log.Println(err)
    }

    c.InfluxDB.WriteMetrics(metrics)

    w.Write([]byte(""))
}
