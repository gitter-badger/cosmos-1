package router

import (
	"strings"
	"net/http"
    "encoding/json"

    "github.com/cosmos-io/cosmos/context"
	"github.com/cosmos-io/cosmos/util"
)

func GetContainers(c context.CosmosContext,
    w http.ResponseWriter,
    r *http.Request) {
	containers, err := c.CosmosService.GetContainers()
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

func AddContainersOfPlanet(c context.CosmosContext,
    w http.ResponseWriter,
    r *http.Request) {
	r.ParseForm()
    planet := c.Params["planet"]
	body, err := util.GetBodyFromRequest(r)

	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
	}

	err = c.CosmosService.AddContainersOfPlanet(planet, body)
	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
	}

    w.Write([]byte(""))
}

func GetContainersOfPlanet(c context.CosmosContext,
    w http.ResponseWriter,
    r *http.Request) {
    planet := c.Params["planet"]

	containers, err := c.CosmosService.GetContainersOfPlanet(planet, true)
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

func GetContainerMetrics(c context.CosmosContext,
    w http.ResponseWriter,
    r *http.Request) {
	r.ParseForm()

	metric := strings.Split(util.GetQueryParam(r, "metric", "all"), ",")
	period := util.GetQueryParam(r, "period", "10m")

    planet := c.Params["planet"]
	container := strings.Replace(c.Params["container"], ".", "_", -1)

	metrics, err := c.CosmosService.GetContainerMetrics(planet, container, metric, period)
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
