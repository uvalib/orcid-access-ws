package test

import (
	"net/http"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/client"
	"testing"
)

//
// delete ORCID attributes tests
//

func TestDeleteOrcidAtributesHappyDay(t *testing.T) {

	expected := http.StatusOK

	id := randomCid()
	attributes := randomOrcidAttributes()
	status := client.SetOrcidAttributes(cfg.Endpoint, id, goodToken, attributes)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	status = client.DelOrcidAttributes(cfg.Endpoint, id, goodToken)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestDeleteOrcidAttributesEmptyId(t *testing.T) {
	expected := http.StatusBadRequest
	status := client.DelOrcidAttributes(cfg.Endpoint, empty, goodToken)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//func TestDeleteOrcidAttributesNotFoundId( t *testing.T ) {
//    expected := http.StatusNotFound
//    status := client.DelOrcidAttributes( cfg.Endpoint, badCid, goodToken )
//    if status != expected {
//        t.Fatalf( "Expected %v, got %v\n", expected, status )
//    }
//}

func TestDeleteOrcidAttributesEmptyToken(t *testing.T) {
	expected := http.StatusBadRequest
	status := client.DelOrcidAttributes(cfg.Endpoint, goodCid, empty)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestDeleteOrcidAttributesBadToken(t *testing.T) {
	expected := http.StatusForbidden
	status := client.DelOrcidAttributes(cfg.Endpoint, goodCid, badToken)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//
// end of file
//
