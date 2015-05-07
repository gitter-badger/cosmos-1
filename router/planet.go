package router

import (
    "strconv"
	"net/http"
    "encoding/json"

	"github.com/cosmos-io/cosmos/context"
)

func GetPlanets(
    c context.Context,
    w http.ResponseWriter,
    r *http.Request) {
    planets, err := c.InfluxDB.QueryPlanets()
    if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        contentLength := strconv.Itoa(len(js))
        w.WriteHeader(http.StatusInternalServerError)
        w.Header().Set("Content-Length", contentLength)
        w.Write(js)
		return
    }

    res := map[string]interface{} {
        "data": planets,
    }

    js, _ := json.Marshal(res)
    contentLength := strconv.Itoa(len(js))
    w.Header().Set("Content-Length", contentLength)
    w.Write(js)
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
