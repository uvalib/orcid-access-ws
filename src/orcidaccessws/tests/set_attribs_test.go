package test

import (
   "net/http"
   "orcidaccessws/api"
   "orcidaccessws/client"
   "testing"
)

//
// set ORCID attributes tests
//

func TestSetOrcidAttributesNew(t *testing.T) {
   expected := http.StatusOK
   id := randomCid()
   attributes := randomOrcidAttributes()

   status := client.SetOrcidAttributes(cfg.Endpoint, id, goodToken, attributes )
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   status, current := client.GetOrcidAttributes(cfg.Endpoint, id, goodToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   if current == nil || len(current) == 0 {
      t.Fatalf("Expected to find orcid for %s and did not\n", attributes.Id)
   }

   ensureIdenticalOrcidsAttributes( t, current[ 0 ], &attributes )
}

func TestSetOrcidAttributesUpdate(t *testing.T) {
   expected := http.StatusOK
   cid := randomCid()
   attributes1 := randomOrcidAttributes()
   attributes2 := randomOrcidAttributes()

   status := client.SetOrcidAttributes(cfg.Endpoint, cid, goodToken, attributes1 )
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   status, current := client.GetOrcidAttributes(cfg.Endpoint, cid, goodToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   if current == nil || len(current) == 0 {
      t.Fatalf("Expected to find orcid for %s and did not\n", attributes1.Id)
   }

   ensureIdenticalOrcidsAttributes( t, current[ 0 ], &attributes1 )

   status = client.SetOrcidAttributes(cfg.Endpoint, cid, goodToken, attributes2 )
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   status, current = client.GetOrcidAttributes(cfg.Endpoint, cid, goodToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   if current == nil || len(current) == 0 {
      t.Fatalf("Expected to find orcid for %s and did not\n", attributes2.Id)
   }

   ensureIdenticalOrcidsAttributes( t, current[ 0 ], &attributes2 )
}

func TestSetOrcidAttributesEmptyId(t *testing.T) {
   expected := http.StatusBadRequest
   attributes := randomOrcidAttributes()
   status := client.SetOrcidAttributes(cfg.Endpoint, empty, goodToken, attributes )
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

func TestSetOrcidAttributesEmptyOrcid(t *testing.T) {
   expected := http.StatusBadRequest
   attributes := api.OrcidAttributes{ Orcid: empty }
   status := client.SetOrcidAttributes(cfg.Endpoint, goodCid, goodToken, attributes )
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

func TestSetOrcidAttributesEmptyToken(t *testing.T) {
   expected := http.StatusBadRequest
   attributes := randomOrcidAttributes()
   status := client.SetOrcidAttributes(cfg.Endpoint, goodCid, empty, attributes )
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

func TestSetOrcidAttributesBadToken(t *testing.T) {
   expected := http.StatusForbidden
   attributes := randomOrcidAttributes()
   status := client.SetOrcidAttributes(cfg.Endpoint, goodCid, badToken, attributes )
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

//
// end of file
//