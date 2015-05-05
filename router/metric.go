package router

import (
    "log"
    "net/http"
    "encoding/json"
    
    "github.com/cosmos-io/cosmos/context"
)

type Metric struct {
    Planet string
    Container string
}

func PostMetric(
    c context.CosmosContext,
    w http.ResponseWriter,
    r *http.Request) {

    var metric *Metric
    err := json.Unmarshal(c.Body, &metric)
    if err != nil {
        log.Println(err)
    }

    log.Println(metric)

    w.Write([]byte(""))
}
