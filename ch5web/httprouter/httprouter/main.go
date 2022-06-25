package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func hello(wr http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	wr.Write([]byte("hello"))
}

func main() {
	r := httprouter.New()
	r.PUT("/user/installations/:installation_id/repositories/:reposit", hello)
	r.GET("/marketplace_listing/plans/", hello)
	r.GET("/search", hello)
	r.GET("/status", hello)
	r.GET("/support", hello)

	http.ListenAndServe(":8080", r)
}
