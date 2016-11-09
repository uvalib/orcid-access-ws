package main

import (
    "io/ioutil"
    "log"
    "testing"
    "orcidaccessws/client"
    "gopkg.in/yaml.v2"
    "net/http"
    "strings"
    "orcidaccessws/api"
    //"fmt"
)

type TestConfig struct {
    Endpoint  string
    Token     string
}

var cfg = loadConfig( )

var goodCid = "dpg3k"
var badCid = "badness"
var goodToken = cfg.Token
var badToken = "badness"
var empty = " "

//
// healthcheck tests
//

func TestHealthCheck( t *testing.T ) {
    expected := http.StatusOK
    status := client.HealthCheck( cfg.Endpoint )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

//
// version tests
//

func TestVersionCheck( t *testing.T ) {
    expected := http.StatusOK
    status, version := client.VersionCheck( cfg.Endpoint )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }

    if len( version ) == 0 {
        t.Fatalf( "Expected non-zero length version string\n" )
    }
}

//
// statistics tests
//

func TestStatistics( t *testing.T ) {
    expected := http.StatusOK
    status, _ := client.Statistics( cfg.Endpoint )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }

    //if len( version ) == 0 {
    //    t.Fatalf( "Expected non-zero length version string\n" )
    //}
}

//
// get single ORCID tests
//

func TestGetOrcidHappyDay( t *testing.T ) {

    expected := http.StatusOK
    id := goodCid
    status, orcids := client.GetOneOrcid( cfg.Endpoint, id, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }

    if orcids == nil {
        t.Fatalf( "Expected to find orcid for %s and did not\n", id )
    }

    ensureValidOrcids( t, orcids )
}

func TestGetOrcidEmptyId( t *testing.T ) {
    expected := http.StatusBadRequest
    status, _ := client.GetOneOrcid( cfg.Endpoint, empty, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestGetOrcidNotFoundId( t *testing.T ) {
    expected := http.StatusNotFound
    status, _ := client.GetOneOrcid( cfg.Endpoint, badCid, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestGetOrcidEmptyToken( t *testing.T ) {
    expected := http.StatusBadRequest
    status, _ := client.GetOneOrcid( cfg.Endpoint, goodCid, empty )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestGetOrcidBadToken( t *testing.T ) {
    expected := http.StatusForbidden
    status, _ := client.GetOneOrcid( cfg.Endpoint, goodCid, badToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

//
// get all ORCID tests
//

func TestGetAllOrcidHappyDay( t *testing.T ) {

    expected := http.StatusOK
    status, orcids := client.GetAllOrcid( cfg.Endpoint, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }

    ensureValidOrcids( t, orcids )
}

func TestGetAllOrcidEmptyToken( t *testing.T ) {
    expected := http.StatusBadRequest
    status, _ := client.GetAllOrcid( cfg.Endpoint, empty )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestGetAllOrcidBadToken( t *testing.T ) {
    expected := http.StatusForbidden
    status, _ := client.GetAllOrcid( cfg.Endpoint, badToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

/*

//
// DOI create tests
//

func TestCreateHappyDay( t *testing.T ) {
    expected := http.StatusOK
    status, entity := client.Create( cfg.Endpoint, goodShoulder, goodToken)
    if status != expected {
        t.Fatalf("Expected %v, got %v\n", expected, status)
    }

    if entity == nil {
        t.Fatalf("Expected to create entity successfully and did not\n" )
    }

    if emptyField( entity.Id ) {
        t.Fatalf( "Expected non-empty ID field but it is empty\n" )
    }
}

func TestCreateEmptyToken( t *testing.T ) {
    expected := http.StatusBadRequest
    status, _ := client.Create( cfg.Endpoint, goodShoulder, empty )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestCreateBadToken( t *testing.T ) {
    expected := http.StatusForbidden
    status, _ := client.Create( cfg.Endpoint, goodShoulder, badToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

//
// DOI update tests
//

func TestUpdateHappyDay( t *testing.T ) {

    doi := createGoodDoi( t )
    entity := testEntity( )
    entity.Id = doi

    expected := http.StatusOK
    status := client.Update( cfg.Endpoint, entity, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestUpdateEmptyId( t *testing.T ) {
    entity := testEntity( )
    entity.Id = empty
    expected := http.StatusBadRequest
    status := client.Update( cfg.Endpoint, entity, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestUpdateBadId( t *testing.T ) {
    expected := http.StatusBadRequest
    status := client.Update( cfg.Endpoint, api.Entity{ Id: badDoi }, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestUpdateEmptyToken( t *testing.T ) {
    entity := testEntity( )
    entity.Id = plausableDoi
    expected := http.StatusBadRequest
    status := client.Update( cfg.Endpoint, entity, empty )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestUpdateBadToken( t *testing.T ) {
    entity := testEntity( )
    entity.Id = plausableDoi
    expected := http.StatusForbidden
    status := client.Update( cfg.Endpoint, entity, badToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

//
// DOI delete tests
//

func TestDeleteHappyDay( t *testing.T ) {
    expected := http.StatusOK
    doi := createGoodDoi( t )
    status := client.Delete( cfg.Endpoint, doi, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestDeleteEmptyId( t *testing.T ) {
    expected := http.StatusBadRequest
    status := client.Delete( cfg.Endpoint, empty, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestDeleteBadId( t *testing.T ) {
    expected := http.StatusBadRequest
    status := client.Delete( cfg.Endpoint, badDoi, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestDeleteEmptyToken( t *testing.T ) {
    expected := http.StatusBadRequest
    status := client.Delete( cfg.Endpoint, plausableDoi, empty )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestDeleteBadToken( t *testing.T ) {
    expected := http.StatusForbidden
    status := client.Delete( cfg.Endpoint, plausableDoi, badToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

//
// DOI revoke tests
//

func TestRevokeHappyDay( t *testing.T ) {

    expected := http.StatusOK
    doi := createGoodDoi( t )
    entity := testEntity( )
    entity.Id = doi

    status := client.Update( cfg.Endpoint, entity, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }

    status = client.Revoke( cfg.Endpoint, entity.Id, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestRevokeEmptyId( t *testing.T ) {
    expected := http.StatusBadRequest
    status := client.Revoke( cfg.Endpoint, empty, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestRevokeBadId( t *testing.T ) {
    expected := http.StatusBadRequest
    status := client.Revoke( cfg.Endpoint, badDoi, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestRevokeEmptyToken( t *testing.T ) {
    expected := http.StatusBadRequest
    status := client.Revoke( cfg.Endpoint, plausableDoi, empty )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestRevokeBadToken( t *testing.T ) {
    expected := http.StatusForbidden
    status := client.Revoke( cfg.Endpoint, plausableDoi, badToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

*/

//
// helpers
//

func ensureValidOrcids( t *testing.T, orcids [] * api.Orcid ) {

    for _, e := range orcids {
        if emptyField( e.Id ) ||
           emptyField( e.Cid ) ||
           emptyField( e.Orcid ) ||
           emptyField( e.CreatedAt ) {
           log.Printf( "%t", e )
           t.Fatalf( "Expected non-empty field but one is empty\n" )
        }
    }
}

func emptyField( field string ) bool {
    return len( strings.TrimSpace( field ) ) == 0
}

func loadConfig( ) TestConfig {

    data, err := ioutil.ReadFile( "service_test.yml" )
    if err != nil {
        log.Fatal( err )
    }

    var c TestConfig
    if err := yaml.Unmarshal( data, &c ); err != nil {
        log.Fatal( err )
    }

    log.Printf( "endpoint [%s]\n", c.Endpoint )
    log.Printf( "token    [%s]\n", c.Token )

    return c
}