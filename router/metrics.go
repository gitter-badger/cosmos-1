package router

import (
    "net/http"
    "encoding/json"
    
    "github.com/cosmos-io/cosmos/context"
    "github.com/cosmos-io/cosmos/model"
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
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
    }

    c.InfluxDB.WriteMetrics(metrics)

    w.Write([]byte(""))
}
