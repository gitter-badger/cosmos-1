package router

import (
    "strconv"
	"net/http"
    "encoding/json"

    "github.com/cosmos-io/cosmos/context"
)

func GetContainers(
    c context.Context,
    w http.ResponseWriter,
    r *http.Request) {
    planet := c.GetQueryParam("planet", "")
    
    if planet == "" {
        res := map[string]string { "error": "" }
        js, _ := json.Marshal(res)
        contentLength := strconv.Itoa(len(js))
        w.WriteHeader(http.StatusInternalServerError)
        w.Header().Set("Content-Length", contentLength)
        w.Write(js)
		return
    }

    containers, err := c.InfluxDB.QueryContainers(planet)
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
        "data": containers,
    }

    js, _ := json.Marshal(res)
    contentLength := strconv.Itoa(len(js))
    w.Header().Set("Content-Length", contentLength)
    w.Write(js)
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
