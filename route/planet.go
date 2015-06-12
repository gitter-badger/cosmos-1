package route

import (
	"net/http"

	"github.com/cosmos-io/cosmos/context"
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

// legacy
/*func GetPlanets(c context.CosmosContext,
    w http.ResponseWriter,
    r *http.Request) {
	planets, err := c.CosmosService.GetPlanets()
	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.Write(js)
		return
	}

    js, err := json.Marshal(planets)
    if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.Write(js)
		return
    }

    w.Write(js)
}*/

/*func GetPlanetMetrics(
    c context.Context,
    w http.ResponseWriter,
    r *http.Request) {
    planet := c.Params["planet"]
	metric := strings.Split(c.GetQueryParam("metric", "all"), ",")

	metrics, err := c.CosmosService.GetPlanetMetrics(planet, metric)
	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
	}

    js, err := json.Marshal(metrics)
    if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
    }

    w.Write(js)
}*/
