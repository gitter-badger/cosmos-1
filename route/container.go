package route

import (
	"net/http"

    "github.com/cosmoshq/cosmos/context"
)

func GetContainers(
    c context.Context,
    w http.ResponseWriter,
    r *http.Request) (int, map[string]interface{}) {
    planet := c.GetQueryParam("planet", "")

    var status int
    var res map[string]interface{}
    
    if planet == "" {
        status = http.StatusBadRequest
        res := map[string]interface{} { "error": "planet is empty." }
        return status, res
    }

    containers, err := c.InfluxDB.QueryContainers(planet)
    if err != nil {
        status = http.StatusBadRequest
        res = map[string]interface{} { "error": err.Error() }
        return status, res
    }

    status = http.StatusOK
    res = map[string]interface{} {
        "data": containers,
    }

    return status, res
}
