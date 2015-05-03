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

func GetContainers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    cosmos := context.Get(r, "cosmos").(*service.CosmosService)
    
	containers, err := cosmos.GetContainers()
	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
	}

    js, err := json.Marshal(containers)
    if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
    }
    
    w.Write(js)
}

func AddContainersOfPlanet(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
	r.ParseForm()
    planet := mux.Vars(r)["planet"]
	body, err := util.GetBodyFromRequest(r)

	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
	}

    cosmos := context.Get(r, "cosmos").(*service.CosmosService)
	err = cosmos.AddContainersOfPlanet(planet, body)
	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
	}

    w.Write([]byte(""))
}

func GetContainersOfPlanet(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    planet := mux.Vars(r)["planet"]

    cosmos := context.Get(r, "cosmos").(*service.CosmosService)
	containers, err := cosmos.GetContainersOfPlanet(planet, true)
	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
	}

    js, err := json.Marshal(containers)
    if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
    }

    w.Write(js)
}

func GetContainerMetrics(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
	r.ParseForm()

	metric := strings.Split(util.GetQueryParam(r, "metric", "all"), ",")
	period := util.GetQueryParam(r, "period", "10m")

    planet := mux.Vars(r)["planet"]
	container := strings.Replace(mux.Vars(r)["container"], ".", "_", -1)

    cosmos := context.Get(r, "cosmos").(*service.CosmosService)
	metrics, err := cosmos.GetContainerMetrics(planet, container, metric, period)
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
}
