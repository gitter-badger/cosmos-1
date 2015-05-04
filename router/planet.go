package router

import (
	"strings"
	"net/http"
    "encoding/json"

	"github.com/cosmos-io/cosmos/context"
	"github.com/cosmos-io/cosmos/util"   
)

func GetPlanets(c context.CosmosContext,
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
}

func GetPlanetMetrics(c context.CosmosContext,
    w http.ResponseWriter,
    r *http.Request) {
    planet := c.Params["planet"]
	metric := strings.Split(util.GetQueryParam(r, "metric", "all"), ",")

	metrics, err := c.CosmosService.GetPlanetMetrics(planet, metric)
	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.Write(js)
		return
	}

    js, err := json.Marshal(metrics)
    if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.Write(js)
		return
    }

    w.Write(js)
}
