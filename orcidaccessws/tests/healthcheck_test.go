package test

import (
	"net/http"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/client"
	"testing"
)

//
// healthcheck tests
//

func TestHealthCheck(t *testing.T) {
	expected := http.StatusOK
	status := client.HealthCheck(cfg.Endpoint)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//
// end of file
//
