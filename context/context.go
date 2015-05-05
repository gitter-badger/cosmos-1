package context

import (
    "github.com/cosmos-io/cosmos/service"
    "github.com/cosmos-io/cosmos/influxdb"
)

type CosmosContext struct {
    CosmosService *service.CosmosService
    InfluxDB *influxdb.InfluxDB
    Params map[string]string
    Body []byte
    QueryParams map[string][]string
}

func (c *CosmosContext) GetQueryParam(key string, defaultValue string) string {
    values := c.QueryParams[key]
    if values == nil {
        return defaultValue
    }
    return values[0]
}
