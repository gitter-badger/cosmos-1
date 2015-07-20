package route

import (
	"net/http"

	"github.com/cosmoshq/cosmos/context"
)

func GetPlanets(
    c context.Context,
    w http.ResponseWriter,
    r *http.Request) (int, map[string]interface{}) {
    planets, err := c.InfluxDB.QueryPlanets()

    var status int
    var res map[string]interface{}
    
    if err != nil {
        status = http.StatusBadRequest
        res = map[string]interface{} { "error": err.Error() }
        return status, res
    }

    status = http.StatusOK    
    res = map[string]interface{} {
        "data": planets,
    }

    return status, res
}
