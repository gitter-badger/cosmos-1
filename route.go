package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

func postContainers(r render.Render, params martini.Params, req *http.Request) {
	req.ParseForm()

	hostId := params["id"]
	json := NewResJson()
	json["host"] = hostId
	r.JSON(200, json)
}

func getContainers(r render.Render, params martini.Params, req *http.Request) {
	req.ParseForm()

	results, err := db.Query("SELECT * FROM hd_used", "")
	if err != nil {
		fmt.Println("Error")
		fmt.Printf("%s\n", err)
		r.JSON(500, err)
		return
	}
	fmt.Println("ok")

	//  hostId := params["id"]
	//  json := NewResJson()
	//	json["host"] = hostId
	r.JSON(200, results)

}
