package router

import (
	//"net/http"
    //"encoding/json"

    //"github.com/cosmos-io/cosmos/context"
)

/*func GetNewsFeeds(c context.Context,
    w http.ResponseWriter,
    r *http.Request) {
    time := ""
	newsfeeds, err := c.CosmosService.GetNewsFeeds(time)
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
}*/
