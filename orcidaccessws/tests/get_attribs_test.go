package test

import (
	"github.com/uvalib/orcid-access-ws/orcidaccessws/client"
	"net/http"
	"testing"
)

//
// get ORCID attribute tests
//

func TestGetOrcidAttributesHappyDay(t *testing.T) {

	expected := http.StatusOK
	id := goodCid
	status, attributes := client.GetOrcidAttributes(cfg.Endpoint, id, goodToken( cfg.Secret ))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	if attributes == nil || len(attributes) == 0 {
		t.Fatalf("Expected to find orcid for %s and did not\n", id)
	}

	ensureValidOrcidsAttributes(t, attributes)
}

func TestGetOrcidAttributesEmptyId(t *testing.T) {
	expected := http.StatusBadRequest
	status, _ := client.GetOrcidAttributes(cfg.Endpoint, empty, goodToken( cfg.Secret ))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestGetOrcidAttributesNotFoundId(t *testing.T) {
	expected := http.StatusNotFound
	status, _ := client.GetOrcidAttributes(cfg.Endpoint, badCid, goodToken( cfg.Secret ))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestGetOrcidAttributesEmptyToken(t *testing.T) {
	expected := http.StatusBadRequest
	status, _ := client.GetOrcidAttributes(cfg.Endpoint, goodCid, empty)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestGetOrcidAttributesBadToken(t *testing.T) {
	expected := http.StatusForbidden
	status, _ := client.GetOrcidAttributes(cfg.Endpoint, goodCid, badToken( cfg.Secret ))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//
// end of file
//
