package router

import (
	"strings"
	"net/http"
    "encoding/json"

	"github.com/cosmos-io/cosmos/service"
	"github.com/cosmos-io/cosmos/util"
    
    "github.com/gorilla/mux"
    "github.com/gorilla/context"
)

func GetPlanets(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    cosmos := context.Get(r, "cosmos").(*service.CosmosService)
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

func GetPlanetMetrics(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    planet := mux.Vars(r)["planet"]
	metric := strings.Split(util.GetQueryParam(r, "metric", "all"), ",")

    cosmos := context.Get(r, "cosmos").(*service.CosmosService)
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
