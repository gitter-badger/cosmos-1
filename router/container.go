package router

import (
	"net/http"

    "github.com/cosmos-io/cosmos/context"
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

// legacy
/*func GetContainers(
    c context.Context,
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
}*/

/*func AddContainersOfPlanet(
    c context.Context,
    w http.ResponseWriter,
    r *http.Request) {
	r.ParseForm()
    planet := c.Params["planet"]
    body := c.Body

	if body == nil {
        res := map[string]string { "error": "HTTP body is invalid." }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
	}

	err := c.CosmosService.AddContainersOfPlanet(planet, body)
	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
	}

    w.Write([]byte(""))
}

func GetContainersOfPlanet(
    c context.Context,
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

func GetContainerMetrics(
    c context.Context,
    w http.ResponseWriter,
    r *http.Request) {
	r.ParseForm()

	metric := strings.Split(c.GetQueryParam("metric", "all"), ",")
	period := c.GetQueryParam("period", "10m")

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
}*/
