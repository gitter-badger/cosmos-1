package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmos-io/cosmos/service"
	"github.com/cosmos-io/cosmos/util"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func GetPlanets(r render.Render, req *http.Request, cosmos *service.CosmosService) {
	token := util.GetQueryParam(req, "token", DEFAULT_USER)

	result, err := cosmos.GetPlanets(token)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, result)
}

func GetPlanetMetrics(r render.Render, params martini.Params, req *http.Request, cosmos *service.CosmosService) {
	token := util.GetQueryParam(req, "token", DEFAULT_USER)
	planet := params["planetName"]

	metric := strings.Split(util.GetQueryParam(req, "metric", "all"), ",")

	result, err := cosmos.GetPlanetMetrics(token, planet, metric)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, result)
}
