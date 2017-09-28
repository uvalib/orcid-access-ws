package test

import (
   "fmt"
   "io/ioutil"
   "log"
   "math/rand"
   "orcidaccessws/api"
   "strconv"
   "strings"
   "testing"
   "time"
    "gopkg.in/yaml.v2"
)

type config struct {
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
// test helpers
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

func loadConfig() config {

   data, err := ioutil.ReadFile("service_test.yml")
   if err != nil {
      log.Fatal(err)
   }

   var c config
   if err := yaml.Unmarshal(data, &c); err != nil {
      log.Fatal(err)
   }

   log.Printf("endpoint [%s]\n", c.Endpoint)
   log.Printf("token    [%s]\n", c.Token)

   return c
}

//
// end of file
//
