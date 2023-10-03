package main

import (
	"encoding/json"
	"fmt"
	"github.com/merlante/prbac-spicedb/api"
	"io"
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

	services, err := getRbacServices()
	if err != nil {
		fmt.Errorf("%v", err)
		os.Exit(1)
	}

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

func getRbacServices() (services server.Services, err error) {
	servicesFile, err := os.Open("services.json")
	if err != nil {
		return nil, err
	}
	defer servicesFile.Close()

	bytes, err := io.ReadAll(servicesFile)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(bytes, &services)

	return
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
