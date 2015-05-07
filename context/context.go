package context

import (
    "github.com/cosmos-io/cosmos/influxdb"
)

type Context struct {
    InfluxDB *influxdb.InfluxDB
    Params map[string]string
    Body []byte
    QueryParams map[string][]string
}

func (c *Context) GetQueryParam(key string, defaultValue string) string {
    values := c.QueryParams[key]
    if values == nil {
        return defaultValue
    }
    return values[0]
}
