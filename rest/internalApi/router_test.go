package internalApi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.wedeliver.com/wedeliver/wallet/utils"
	"gitlab.wedeliver.com/wedeliver/wallet/utils/config"
)

func TestHealthEndpoint(t *testing.T) {

	// Initialise Config with defaults
	cfg, err := config.NewConfig(true)
	if err != nil {
		t.Fatalf("failed to read config : %s", err.Error())
	}
	cfg.LogLevel = "fatal"

	// Inititialse logger
	logger := utils.NewLogger(cfg)

	server := httptest.NewServer(NewInternalAPIRouter(NewInternalAPIHandler(cfg, logger), cfg))
	defer server.Close()

	res, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("failed to request API: %v", err)
	}

	statusCode := res.StatusCode
	if statusCode != 200 {
		t.Fatalf("status code = %v", statusCode)
	}

	var body map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read body bytes : %v", err)
	}

	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		t.Fatalf("failed to unmarshal body bytes : %v", err)
	}

	status := body["status"]
	if status != "running" {
		t.Errorf("wrong status obtained, got %s", status.(string))
	}
}
