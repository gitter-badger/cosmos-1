package router

import (
	"net/http"
    "encoding/json"

	"github.com/cosmos-io/cosmos/service"
)

func GetNewsFeeds(w http.ResponseWriter, req *http.Request, cosmos *service.CosmosService) {
    w.Header().Set("Content-Type", "application/json")
    
	newsfeeds, err := cosmos.GetNewsFeeds("")
	if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
	}

    js, err := json.Marshal(newsfeeds)
    if err != nil {
        res := map[string]string { "error": err.Error() }
        js, _ := json.Marshal(res)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(js)
		return
    }

    w.Write(js)
}
