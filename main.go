package main

import (
	"fmt"
	"github.com/merlante/prbac-spicedb/api"
	"net/http"
	"os"

	"github.com/merlante/prbac-spicedb/server"
)

var (
	spiceDBURL   = "localhost:50051"
	spiceDBToken = "foobar"
)

func main() {
	overwriteVarsFromEnv()

	services := getRbacServices()
	spiceDbClient, err := server.GetSpiceDbClient(spiceDBURL, spiceDBToken)
	if err != nil {
		fmt.Errorf("%v", err)
		os.Exit(1)
	}

	server := server.PrbacSpicedbServer{
		RbacServices:  services,
		SpicedbClient: spiceDbClient,
	}
	r := api.Handler(api.NewStrictHandler(server, nil))

	http.ListenAndServe(":8080", r)
}

func getRbacServices() server.Services {
	services := server.Services{}

	pbFilter := server.Filter{
		Name:         "service",
		Operator:     "equals",
		ResourceType: "dispatcher/service",
		Verb:         "view",
	}

	pbResourcePerm := server.ResourcePerm{
		Permission: "dispatcher_view_runs",
		Filter:     pbFilter,
	}

	pbPermission := server.Permission{}
	pbPermission["run:read"] = pbResourcePerm

	services["playbook-dispatcher"] = pbPermission

	return services
}

func overwriteVarsFromEnv() {
	envSpicedbUrl := os.Getenv("SPICEDB_URL")
	if envSpicedbUrl != "" {
		spiceDBURL = envSpicedbUrl
	}
	envSpicedbPsk := os.Getenv("SPICEDB_PSK")
	if envSpicedbPsk != "" {
		spiceDBToken = envSpicedbPsk
	}
}
