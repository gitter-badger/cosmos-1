package route

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/cosmoshq/cosmos/context"
	"github.com/cosmoshq/cosmos/model"
)

var (
	types = map[string]bool{
		"cpu":    true,
		"memory": true,
	}
)

func PostMetrics(
	c context.Context,
	w http.ResponseWriter,
	r *http.Request) (int, map[string]interface{}) {
	var metrics *model.MetricsParam
	err := json.Unmarshal(c.Body, &metrics)

	var status int
	var res map[string]interface{}

	if err != nil {
		status = http.StatusBadRequest
		res = map[string]interface{}{"error": err.Error()}
		return status, res
	}

	c.InfluxDB.WriteMetrics(metrics)

	status = http.StatusOK

	return status, res
}

func GetMetrics(
	c context.Context,
	w http.ResponseWriter,
	r *http.Request) (int, map[string]interface{}) {
	t := c.GetQueryParam("type", "")
	planet := c.GetQueryParam("planet", "")
	container := c.GetQueryParam("container", "")

	var res map[string]interface{}

	if types[t] == false {
		err := fmt.Sprintf("%s type is not supported.", t)
		res = map[string]interface{}{"error": err}
		return http.StatusBadRequest, res
	}

	if planet == "" {
		err := fmt.Sprintf("planet is empty.")
		res = map[string]interface{}{"error": err}
		return http.StatusBadRequest, res
	}

	var metrics interface{}
	var err error

	if container == "" {
		metrics, err = c.InfluxDB.QueryPlanetMetrics(planet, t)
	} else {
		metrics, err = c.InfluxDB.QueryContainerMetrics(planet, container, t)
	}

	if err != nil {
		res := map[string]interface{}{"error": err.Error()}
		return http.StatusInternalServerError, res
	}

	res = map[string]interface{}{
		"data": metrics,
	}

	return http.StatusOK, res
}
