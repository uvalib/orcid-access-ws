package test

import (
	"net/http"
	//"github.com/uvalib/orcid-access-ws/orcidaccessws/api"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/client"
	"testing"
)

//
// set ORCID attributes tests
//

func TestSetOrcidAttributesNew(t *testing.T) {
	expected := http.StatusOK
	id := randomCid()
	attributes := randomOrcidAttributes()

	status := client.SetOrcidAttributes(cfg.Endpoint, id, goodToken(cfg.Secret), attributes)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	status, current := client.GetOrcidAttributes(cfg.Endpoint, id, goodToken(cfg.Secret))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	if current == nil || len(current) == 0 {
		t.Fatalf("Expected to find orcid for %s and did not\n", attributes.ID)
	}

	ensureIdenticalOrcidsAttributes(t, current[0], &attributes)
}

func TestSetOrcidAttributesUpdate(t *testing.T) {
	expected := http.StatusOK
	cid := randomCid()
	attributes1 := randomOrcidAttributes()
	attributes2 := randomOrcidAttributes()

	status := client.SetOrcidAttributes(cfg.Endpoint, cid, goodToken(cfg.Secret), attributes1)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	status, current := client.GetOrcidAttributes(cfg.Endpoint, cid, goodToken(cfg.Secret))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	if current == nil || len(current) == 0 {
		t.Fatalf("Expected to find orcid for %s and did not\n", attributes1.ID)
	}

	ensureIdenticalOrcidsAttributes(t, current[0], &attributes1)

	status = client.SetOrcidAttributes(cfg.Endpoint, cid, goodToken(cfg.Secret), attributes2)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	status, current = client.GetOrcidAttributes(cfg.Endpoint, cid, goodToken(cfg.Secret))
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	if current == nil || len(current) == 0 {
		t.Fatalf("Expected to find orcid for %s and did not\n", attributes2.ID)
	}

	ensureIdenticalOrcidsAttributes(t, current[0], &attributes2)
}

func TestSetOrcidAttributesEmptyId(t *testing.T) {
	expected := http.StatusBadRequest
	attributes := randomOrcidAttributes()
	status := client.SetOrcidAttributes(cfg.Endpoint, empty, goodToken(cfg.Secret), attributes)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//func TestSetOrcidAttributesEmptyOrcid(t *testing.T) {
//   expected := http.StatusBadRequest
//   attributes := api.OrcidAttributes{ Orcid: empty }
//   status := client.SetOrcidAttributes(cfg.Endpoint, goodCid, goodToken( cfg.Secret ), attributes )
//   if status != expected {
//      t.Fatalf("Expected %v, got %v\n", expected, status)
//   }
//}

func TestSetOrcidAttributesEmptyToken(t *testing.T) {
	expected := http.StatusBadRequest
	attributes := randomOrcidAttributes()
	status := client.SetOrcidAttributes(cfg.Endpoint, goodCid, empty, attributes)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestSetOrcidAttributesBadToken(t *testing.T) {
	expected := http.StatusForbidden
	attributes := randomOrcidAttributes()
	status := client.SetOrcidAttributes(cfg.Endpoint, goodCid, badToken(cfg.Secret), attributes)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//
// end of file
//
