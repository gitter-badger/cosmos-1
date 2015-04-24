package router

import (
	"fmt"
	"net/http"

	"github.com/cosmos-io/cosmos/service"
	"github.com/cosmos-io/cosmos/util"
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
