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

func GetContainers(r render.Render, req *http.Request, cosmos *service.CosmosService) {
	token := util.GetToken(req)

	result, err := cosmos.GetContainers(token)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, result)
}

func AddContainersOfPlanet(r render.Render, params martini.Params, req *http.Request, cosmos *service.CosmosService) {
	req.ParseForm()
	token := util.GetToken(req)
	planet := params["planetName"]

	body, err := util.GetBodyFromRequest(req)

	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	err = cosmos.AddContainersOfPlanet(token, planet, body)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.Status(http.StatusOK)
}

func GetContainersOfPlanet(r render.Render, params martini.Params, req *http.Request, cosmos *service.CosmosService) {
	token := util.GetToken(req)
	planet := params["planetName"]

	result, err := cosmos.GetContainersOfPlanet(token, planet, true)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, result)
}

func GetContainerMetrics(r render.Render, params martini.Params, req *http.Request, cosmos *service.CosmosService) {
	req.ParseForm()

	metric := util.GetMetric(req)
	token := util.GetToken(req)
	planetName := params["planetName"]
	containerName := strings.Replace(params["containerName"], ".", "_", -1)

	result, err := cosmos.GetContainerMetrics(token, planetName, containerName, metric)
	if err != nil {
		fmt.Println(err)
		r.JSON(http.StatusInternalServerError, err)
		return
	}

	r.JSON(http.StatusOK, result)
}
