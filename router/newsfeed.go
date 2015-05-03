package router

import (
	"net/http"
    "encoding/json"

    "github.com/cosmos-io/cosmos/service"
    
    "github.com/gorilla/context"
)

func GetNewsFeeds(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    cosmos := context.Get(r, "cosmos").(*service.CosmosService)
    time := ""
	newsfeeds, err := cosmos.GetNewsFeeds(time)
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
