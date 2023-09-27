package main

import (
	"github.com/merlante/prbac-spicedb/api"
	"net/http"
)

func main() {
	r := api.Handler(api.Unimplemented{})
	http.ListenAndServe(":8080", r)
}
