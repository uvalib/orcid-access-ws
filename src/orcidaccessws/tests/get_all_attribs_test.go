package test

import (
	"net/http"
	"orcidaccessws/client"
	"testing"
)

//
// get all ORCID attributes tests
//

func TestGetAllOrcidAttributesHappyDay(t *testing.T) {

	expected := http.StatusOK
	status, attributes := client.GetAllOrcidAttributes(cfg.Endpoint, goodToken)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	ensureValidOrcidsAttributes(t, attributes)
}

func TestGetAllOrcidAttributesEmptyToken(t *testing.T) {
	expected := http.StatusBadRequest
	status, _ := client.GetAllOrcidAttributes(cfg.Endpoint, empty)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestGetAllOrcidAttributesBadToken(t *testing.T) {
	expected := http.StatusForbidden
	status, _ := client.GetAllOrcidAttributes(cfg.Endpoint, badToken)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//
// end of file
//
