package test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

type config struct {
	Endpoint string
	Secret   string
}

var cfg = loadConfig()

var goodCid = "dpg3k"
var badCid = "badness"
var goodOrcid = "0000-0003-4520-4923"
var badOrcid = "9999-9999-0000-0000"
var goodSearch = "Ellen Ramsey"
var notFoundSearch = "hurunglyzit"
var empty = " "
var goodSearchStart = "0"
var badSearchStart = "x"
var goodSearchMax = "25"
var badSearchMax = "x"

var lcAlphaChars = "abcdefghijklmnopqrstuvwxyz"
var numericChars = "0123456789"
var numericAndLcAlphaChars = numericChars + lcAlphaChars

//
// test helpers
//

func randomOrcidAttributes() api.OrcidAttributes {

	// see the RNG
	rand.Seed(time.Now().UnixNano())

	// list of possible characters for the ORCID
	possible := []rune(numericChars)
	orcid := fmt.Sprintf("%s-%s-%s-%s", randomString(possible, 4), randomString(possible, 4),
		randomString(possible, 4), randomString(possible, 4))

	// list of possible characters for the access tokens
	possible = []rune(numericAndLcAlphaChars)
	oauthAccess := randomString(possible, 32)
	oauthRenew := randomString(possible, 32)
	oauthScope := "/my/scope"
	return api.OrcidAttributes{Orcid: orcid,
		OauthAccessToken: oauthAccess, OauthRefreshToken: oauthRenew, OauthScope: oauthScope}
}

func randomCid() string {

	// see the RNG
	rand.Seed(time.Now().UnixNano())

	// list of possible characters
	possible := []rune(numericAndLcAlphaChars)

	return randomString(possible, 5)
}

func randomUpdateCode() string {

	// see the RNG
	rand.Seed(time.Now().UnixNano())

	// list of possible characters
	possible := []rune(numericChars)

	return randomString(possible, 6)
}

func workActivity() api.ActivityUpdate {

	// see the RNG
	rand.Seed(time.Now().UnixNano())

	possible := []rune(lcAlphaChars)
	title := fmt.Sprintf("Title-%s", randomString(possible, 32))
	abstract := fmt.Sprintf("Abstract-%s", randomString(possible, 32))
	pubDate := "2017-03-05"
	url := fmt.Sprintf("www.foobar.com/%s", randomString(possible, 16))
	persons := makePeople(2)
	rt := "journal-article"

	work := api.WorkSchema{Title: title, Abstract: abstract, PublicationDate: pubDate, URL: url, Authors: persons, ResourceType: rt}
	return api.ActivityUpdate{Work: work}
}

func makePeople(number int) []api.Person {

	people := make([]api.Person, number, number)
	for ix := 0; ix < number; ix++ {

		p := api.Person{
			Index:     ix,
			FirstName: fmt.Sprintf("first-%d", ix+1),
			LastName:  fmt.Sprintf("last-%d", ix+1),
		}
		people[ix] = p
	}
	return people
}

func randomString(possible []rune, sz int) string {

	b := make([]rune, sz)
	for i := range b {
		b[i] = possible[rand.Intn(len(possible))]
	}
	return string(b)
}

func ensureIdenticalOrcidsAttributes(t *testing.T, attributes1 *api.OrcidAttributes, attributes2 *api.OrcidAttributes) {

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
		if emptyField(e.ID) ||
			emptyField(e.Cid) ||
			emptyField(e.Orcid) ||
			emptyField(e.URI) ||
			emptyField(e.CreatedAt) {
			log.Printf("%v", e)
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
		emptyField(orcid.URI) ||
		emptyField(orcid.DisplayName) ||
		emptyField(orcid.FirstName) ||
		emptyField(orcid.LastName) {
		log.Printf("%v", orcid)
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

	log.Printf("endpoint  [%s]\n", c.Endpoint)
	log.Printf("secret    [%s]\n", c.Secret)
	return c
}

func badToken(secret string) string {

	// Declare the expiration time of the token
	expirationTime := time.Now().Add(-5 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &jwt.StandardClaims{
		// In JWT, the expiry time is expressed as unix milliseconds
		ExpiresAt: expirationTime.Unix(),
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}

func goodToken(secret string) string {

	// Declare the expiration time of the token
	expirationTime := time.Now().Add(5 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &jwt.StandardClaims{
		// In JWT, the expiry time is expressed as unix milliseconds
		ExpiresAt: expirationTime.Unix(),
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}

//
// end of file
//
