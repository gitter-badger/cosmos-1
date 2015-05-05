package context

import (
    "github.com/cosmos-io/cosmos/service"
    "github.com/cosmos-io/cosmos/influxdb"
)

type CosmosContext struct {
    CosmosService *service.CosmosService
    Params map[string]string
    InfluxDB *influxdb.InfluxDB
}
