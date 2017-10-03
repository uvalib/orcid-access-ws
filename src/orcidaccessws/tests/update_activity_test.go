package test

import (
	"net/http"
	"orcidaccessws/client"
	"testing"
)

//
// update activity tests
//

func TestUpdateActivityNew(t *testing.T) {

	expected := http.StatusOK
	id := goodCid
	newActivity := workActivity()
	status, code := client.UpdateActivity(cfg.Endpoint, id, goodToken, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	if len(code) == 0 {
		t.Fatalf("Expected to receive an update code and did not\n")
	}
}

func TestUpdateActivityUpdate(t *testing.T) {

	expected := http.StatusOK
	id := goodCid
	newActivity := workActivity()
	status, code := client.UpdateActivity(cfg.Endpoint, id, goodToken, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	if len(code) == 0 {
		t.Fatalf("Expected to receive an update code and did not\n")
	}

	newActivity.UpdateCode = code
	status, code = client.UpdateActivity(cfg.Endpoint, id, goodToken, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}

	if len(code) == 0 {
		t.Fatalf("Expected to receive an update code and did not\n")
	}

	if code != newActivity.UpdateCode {
		t.Fatalf("Unexpected update code; was %s should be %s\n", newActivity.UpdateCode, code)
	}
}

func TestUpdateActivityEmptyWorkTitle(t *testing.T) {
	expected := http.StatusBadRequest
	id := goodCid
	newActivity := workActivity()
	newActivity.Work.Title = ""
	status, _ := client.UpdateActivity(cfg.Endpoint, id, goodToken, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestUpdateActivityEmptyWorkResourceType(t *testing.T) {
	expected := http.StatusBadRequest
	id := goodCid
	newActivity := workActivity()
	newActivity.Work.ResourceType = ""
	status, _ := client.UpdateActivity(cfg.Endpoint, id, goodToken, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestUpdateActivityEmptyWorkUrl(t *testing.T) {
	expected := http.StatusBadRequest
	id := goodCid
	newActivity := workActivity()
	newActivity.Work.Url = ""
	status, _ := client.UpdateActivity(cfg.Endpoint, id, goodToken, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestUpdateActivityEmptyId(t *testing.T) {
	expected := http.StatusBadRequest
	newActivity := workActivity()
	status, _ := client.UpdateActivity(cfg.Endpoint, empty, goodToken, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestUpdateActivityNotFoundId(t *testing.T) {
	expected := http.StatusNotFound
	newActivity := workActivity()
	status, _ := client.UpdateActivity(cfg.Endpoint, badCid, goodToken, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestUpdateActivityEmptyToken(t *testing.T) {
	expected := http.StatusBadRequest
	newActivity := workActivity()
	status, _ := client.UpdateActivity(cfg.Endpoint, goodCid, empty, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestUpdateActivityBadToken(t *testing.T) {
	expected := http.StatusForbidden
	newActivity := workActivity()
	status, _ := client.UpdateActivity(cfg.Endpoint, goodCid, badToken, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

func TestUpdateActivityBadWorkResourceType(t *testing.T) {
	expected := http.StatusBadRequest
	id := goodCid
	newActivity := workActivity()
	newActivity.Work.ResourceType = "a-bad-resource-type"
	status, _ := client.UpdateActivity(cfg.Endpoint, id, goodToken, newActivity)
	if status != expected {
		t.Fatalf("Expected %v, got %v\n", expected, status)
	}
}

//
// end of file
//
