package router

import (
	"strings"    
	"net/http"
    "encoding/json"

	"github.com/cosmos-io/cosmos/service"
	"github.com/cosmos-io/cosmos/util"
    
	"github.com/go-martini/martini"
)

func GetPlanets(w http.ResponseWriter, r *http.Request, cosmos *service.CosmosService) {
    w.Header().Set("Content-Type", "application/json")

	planets, err := cosmos.GetPlanets()
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

func GetPlanetMetrics(w http.ResponseWriter, params martini.Params, r *http.Request, cosmos *service.CosmosService) {
    w.Header().Set("Content-Type", "application/json")
    
	planet := params["planetName"]
	metric := strings.Split(util.GetQueryParam(r, "metric", "all"), ",")

	metrics, err := cosmos.GetPlanetMetrics(planet, metric)
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
