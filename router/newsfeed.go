package router

import (
	"net/http"

	"github.com/cosmos-io/cosmos/service"
	"github.com/cosmos-io/cosmos/util"
	"github.com/martini-contrib/render"
)

func GetNewsFeeds(r render.Render, req *http.Request, cosmos *service.CosmosService) {
	token := util.GetToken(req)
	result, err := cosmos.GetNewsFeeds(token, "")
	if err != nil {
		r.JSON(http.StatusInternalServerError, JSON{"error": err.Error()})
		return
	}

	r.JSON(http.StatusOK, result)
}
