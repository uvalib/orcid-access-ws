package test

import (
	"net/http"
	"orcidaccessws/client"
	"testing"
)

//
// search ORCID tests
//

//func TestSearchOrcidHappyDay(t *testing.T) {
//
//	expected := http.StatusOK
//	status, orcids, total := client.SearchOrcid(cfg.Endpoint, goodSearch, goodSearchStart, goodSearchMax, goodToken)
//	if status != expected {
//		t.Fatalf("Expected %v, got %v\n", expected, status)
//	}
//
//	ensureValidSearchResults(t, orcids, goodSearchMax, total)
//}

//func TestSearchOrcidMaxRows(t *testing.T) {
//
//	expected := http.StatusOK
//	status, orcids, total := client.SearchOrcid(cfg.Endpoint, goodSearch, goodSearchStart, goodSearchMax, goodToken)
//	if status != expected {
//		t.Fatalf("Expected %v, got %v\n", expected, status)
//	}
//
//	ensureValidSearchResults(t, orcids, goodSearchMax, total)
//}

func TestSearchOrcidBadStart(t *testing.T) {
	expected := http.StatusBadRequest
	status, _, _ := client.SearchOrcid(cfg.Endpoint, goodSearch, badSearchStart, goodSearchMax, empty)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestSearchOrcidBadMax(t *testing.T) {
	expected := http.StatusBadRequest
	status, _, _ := client.SearchOrcid(cfg.Endpoint, goodSearch, goodSearchStart, badSearchMax, empty)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestSearchOrcidEmptySearch(t *testing.T) {
	expected := http.StatusBadRequest
	status, _, _ := client.SearchOrcid(cfg.Endpoint, empty, goodSearchStart, goodSearchMax, goodToken)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//func TestSearchOrcidNotFoundSearch(t *testing.T) {
//	expected := http.StatusNotFound
//	status, _, _ := client.SearchOrcid(cfg.Endpoint, notFoundSearch, goodSearchStart, goodSearchMax, goodToken)
//	if status != expected {
//		t.Fatalf("Expected %v, got %v\n", expected, status)
//	}
//}

func TestSearchOrcidEmptyToken(t *testing.T) {
	expected := http.StatusBadRequest
	status, _, _ := client.SearchOrcid(cfg.Endpoint, goodSearch, goodSearchStart, goodSearchMax, empty)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestSearchOrcidBadToken(t *testing.T) {
	expected := http.StatusForbidden
	status, _, _ := client.SearchOrcid(cfg.Endpoint, goodSearch, goodSearchStart, goodSearchMax, badToken)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//
// end of file
//
