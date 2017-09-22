package main

import (
   "fmt"
   "io/ioutil"
   "log"
   "math/rand"
   "net/http"
   "orcidaccessws/api"
   "orcidaccessws/client"
   "strconv"
   "strings"
   "testing"
   "time"
    "gopkg.in/yaml.v2"
)

type TestConfig struct {
   Endpoint string
   Token    string
}

var cfg = loadConfig()

var goodCid = "dpg3k"
var badCid = "badness"
var goodToken = cfg.Token
var badToken = "badness"
var goodOrcid = "0000-0002-0566-4186"
var badOrcid = "9999-9999-0000-0000"
var goodSearch = "Dave Goldstein"
var notFoundSearch = "hurunglyzit"
var empty = " "
var goodSearchStart = "0"
var badSearchStart = "x"
var goodSearchMax = "25"
var badSearchMax = "x"

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
// version tests
//

func TestVersionCheck(t *testing.T) {
   expected := http.StatusOK
   status, version := client.VersionCheck(cfg.Endpoint)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   if len(version) == 0 {
      t.Fatalf("Expected non-zero length version string\n")
   }
}

//
// runtime tests
//

func TestRuntimeCheck(t *testing.T) {
   expected := http.StatusOK
   status, runtime := client.RuntimeCheck(cfg.Endpoint)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   if runtime == nil {
      t.Fatalf("Expected non-nil runtime info\n")
   }

   if len( runtime.Version ) == 0 ||
      runtime.AllocatedMemory == 0 ||
      runtime.CpuCount == 0 ||
      runtime.GoRoutineCount == 0 ||
      runtime.ObjectCount == 0 {
      t.Fatalf("Expected non-zero value in runtime info but one is zero\n")
   }
}

//
// statistics tests
//

func TestStatistics(t *testing.T) {
   expected := http.StatusOK
   status, stats := client.Statistics(cfg.Endpoint)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   if stats.RequestCount == 0 {
      t.Fatalf("Expected non-zero request count\n")
   }
}

//
// get ORCID attribute tests
//

func TestGetOrcidAttributesHappyDay(t *testing.T) {

   expected := http.StatusOK
   id := goodCid
   status, attributes := client.GetOrcidAttributes(cfg.Endpoint, id, goodToken)
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
   status, _ := client.GetOrcidAttributes(cfg.Endpoint, empty, goodToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

func TestGetOrcidAttributesNotFoundId(t *testing.T) {
   expected := http.StatusNotFound
   status, _ := client.GetOrcidAttributes(cfg.Endpoint, badCid, goodToken)
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
   status, _ := client.GetOrcidAttributes(cfg.Endpoint, goodCid, badToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

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
// set ORCID attributes tests
//

func TestSetOrcidAttributesNew(t *testing.T) {
   expected := http.StatusOK
   id := randomCid()
   attributes := randomOrcidAttributes()

   status := client.SetOrcidAttributes(cfg.Endpoint, id, attributes, goodToken)
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

   status := client.SetOrcidAttributes(cfg.Endpoint, cid, attributes1, goodToken)
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

   status = client.SetOrcidAttributes(cfg.Endpoint, cid, attributes2, goodToken)
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
   status := client.SetOrcidAttributes(cfg.Endpoint, empty, attributes, goodToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

func TestSetOrcidAttributesEmptyOrcid(t *testing.T) {
   expected := http.StatusBadRequest
   attributes := api.OrcidAttributes{ Orcid: empty }
   status := client.SetOrcidAttributes(cfg.Endpoint, goodCid, attributes, goodToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

func TestSetOrcidAttributesEmptyToken(t *testing.T) {
   expected := http.StatusBadRequest
   attributes := randomOrcidAttributes()
   status := client.SetOrcidAttributes(cfg.Endpoint, goodCid, attributes, empty)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

func TestSetOrcidAttributesBadToken(t *testing.T) {
   expected := http.StatusForbidden
   attributes := randomOrcidAttributes()
   status := client.SetOrcidAttributes(cfg.Endpoint, goodCid, attributes, badToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

//
// delete ORCID attributes tests
//

func TestDeleteOrcidAtributesHappyDay(t *testing.T) {

   expected := http.StatusOK

   id := randomCid()
   attributes := randomOrcidAttributes()
   status := client.SetOrcidAttributes(cfg.Endpoint, id, attributes, goodToken)
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
// get ORCID details tests
//

func TestGetOrcidDetailsHappyDay(t *testing.T) {

   expected := http.StatusOK
   id := goodOrcid
   status, orcid := client.GetOrcidDetails(cfg.Endpoint, id, goodToken)
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
   status, _ := client.GetOrcidDetails(cfg.Endpoint, empty, goodToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

func TestGetOrcidDetailsNotFoundId(t *testing.T) {
   expected := http.StatusNotFound
   status, _ := client.GetOrcidDetails(cfg.Endpoint, badOrcid, goodToken)
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
   status, _ := client.GetOrcidDetails(cfg.Endpoint, goodOrcid, badToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

//
// search ORCID tests
//

func TestSearchOrcidHappyDay(t *testing.T) {

   expected := http.StatusOK
   status, orcids, total := client.SearchOrcid(cfg.Endpoint, goodSearch, goodSearchStart, goodSearchMax, goodToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   ensureValidSearchResults(t, orcids, goodSearchMax, total)
}

func TestSearchOrcidMaxRows(t *testing.T) {

   expected := http.StatusOK
   status, orcids, total := client.SearchOrcid(cfg.Endpoint, goodSearch, goodSearchStart, goodSearchMax, goodToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   ensureValidSearchResults(t, orcids, goodSearchMax, total)
}

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

func TestSearchOrcidNotFoundSearch(t *testing.T) {
   expected := http.StatusNotFound
   status, _, _ := client.SearchOrcid(cfg.Endpoint, notFoundSearch, goodSearchStart, goodSearchMax, goodToken)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }
}

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
// helpers
//

func randomOrcidAttributes() api.OrcidAttributes {

   // see the RNG
   rand.Seed(time.Now().UnixNano())

   // list of possible characters for the ORCID
   possible := []rune("0123456789")
   orcid := fmt.Sprintf("%s-%s-%s-%s", randomString(possible, 4), randomString(possible, 4),
      randomString(possible, 4), randomString(possible, 4))

   // list of possible characters for the access tokens
   possible = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
   oauth_access := randomString(possible, 32 )
   oauth_renew := randomString(possible, 32 )
   oauth_scope := "/my/scope"
   return api.OrcidAttributes{ Orcid: orcid,
      OauthAccessToken: oauth_access, OauthRefreshToken: oauth_renew, OauthScope: oauth_scope }
}

func randomCid() string {

   // see the RNG
   rand.Seed(time.Now().UnixNano())

   // list of possible characters
   possible := []rune("abcdefghijklmnopqrstuvwxyz0123456789")

   return randomString(possible, 5)
}

func randomString(possible []rune, sz int) string {

   b := make([]rune, sz)
   for i := range b {
      b[i] = possible[rand.Intn(len(possible))]
   }
   return string(b)
}

func ensureIdenticalOrcidsAttributes(t *testing.T, attributes1 *api.OrcidAttributes, attributes2 *api.OrcidAttributes ) {

   //log.Printf("%t", attributes1)
   //log.Printf("%t", attributes2)

   if attributes1.Orcid != attributes2.Orcid ||
      attributes1.OauthAccessToken != attributes2.OauthAccessToken ||
      attributes1.OauthRefreshToken != attributes2.OauthRefreshToken ||
      attributes1.OauthScope != attributes2.OauthScope {
      t.Fatalf("Expected identical attributes but they are not\n")
   }
}

func ensureValidOrcidsAttributes(t *testing.T, orcids []*api.OrcidAttributes) {
   for _, e := range orcids {
      if emptyField(e.Id) ||
         emptyField(e.Cid) ||
         emptyField(e.Orcid) ||
         emptyField(e.Uri) ||
         emptyField(e.CreatedAt) {
         log.Printf("%t", e)
         t.Fatalf("Expected non-empty field but one is empty\n")
      }
   }
}

func ensureValidSearchResults(t *testing.T, orcids []*api.OrcidDetails, expectedMax string, totalFound int) {
   for _, e := range orcids {
      ensureValidOrcidDetails(t, e)
   }

   max, _ := strconv.Atoi(expectedMax)
   actualCount := len(orcids)
   if actualCount > max {
      t.Fatalf("Expected %v results, got %v\n", max, actualCount)
   }

   if totalFound < actualCount {
      t.Fatalf("Incorrect search total count, got %v\n", totalFound)
   }
}

func ensureValidOrcidDetails(t *testing.T, orcid *api.OrcidDetails) {
   if emptyField(orcid.Orcid) ||
      emptyField(orcid.Uri) ||
      //emptyField( orcid.DisplayName ) ||
      emptyField(orcid.FirstName) ||
      emptyField(orcid.LastName) {
      log.Printf("%t", orcid)
      t.Fatalf("Expected non-empty field but one is empty\n")
   }
}

func emptyField(field string) bool {
   return len(strings.TrimSpace(field)) == 0
}

func loadConfig() TestConfig {

   data, err := ioutil.ReadFile("service_test.yml")
   if err != nil {
      log.Fatal(err)
   }

   var c TestConfig
   if err := yaml.Unmarshal(data, &c); err != nil {
      log.Fatal(err)
   }

   log.Printf("endpoint [%s]\n", c.Endpoint)
   log.Printf("token    [%s]\n", c.Token)

   return c
}
