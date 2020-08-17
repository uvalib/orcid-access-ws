package test

import (
	"github.com/uvalib/orcid-access-ws/orcidaccessws/client"
	"net/http"
	"testing"
)

//
// get ORCID details tests
//

func TestGetOrcidDetailsHappyDay(t *testing.T) {

	expected := http.StatusOK
	id := goodOrcid
	status, orcid := client.GetOrcidDetails(cfg.Endpoint, id, goodToken(cfg.Secret))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	if orcid == nil {
		t.Fatalf("Expected to find orcid for %s and did not\n", id)
	}

	ensureValidOrcidDetails(t, orcid)
}

func TestGetOrcidDetailsEmptyId(t *testing.T) {
	expected := http.StatusBadRequest
	status, _ := client.GetOrcidDetails(cfg.Endpoint, empty, goodToken(cfg.Secret))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestGetOrcidDetailsNotFoundId(t *testing.T) {
	expected := http.StatusNotFound
	status, _ := client.GetOrcidDetails(cfg.Endpoint, badOrcid, goodToken(cfg.Secret))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestGetOrcidDetailsEmptyToken(t *testing.T) {
	expected := http.StatusBadRequest
	status, _ := client.GetOrcidDetails(cfg.Endpoint, goodOrcid, empty)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestGetOrcidDetailsBadToken(t *testing.T) {
	expected := http.StatusForbidden
	status, _ := client.GetOrcidDetails(cfg.Endpoint, goodOrcid, badToken(cfg.Secret))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//
// end of file
//
