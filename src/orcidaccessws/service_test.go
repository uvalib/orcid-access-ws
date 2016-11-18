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
    "time"
    "math/rand"
    "fmt"
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
var goodOrcid = "0000-0002-0566-4186"
var badOrcid = "9999-9999-0000-0000"
var goodSearch = "Dave Goldstein"
var notFoundSearch = "hurunglyzit"
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

    if orcids == nil || len( orcids ) == 0 {
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

//
// get ORCID details tests
//

func TestGetOrcidDetailsHappyDay( t *testing.T ) {

    expected := http.StatusOK
    id := goodOrcid
    status, orcids := client.GetOrcidDetails( cfg.Endpoint, id, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }

    if orcids == nil || len( orcids ) == 0 {
        t.Fatalf( "Expected to find orcid for %s and did not\n", id )
    }

    ensureValidOrcidDetails( t, orcids )
}

func TestGetOrcidDetailsEmptyId( t *testing.T ) {
    expected := http.StatusBadRequest
    status, _ := client.GetOrcidDetails( cfg.Endpoint, empty, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestGetOrcidDetailsNotFoundId( t *testing.T ) {
    expected := http.StatusNotFound
    status, _ := client.GetOrcidDetails( cfg.Endpoint, badOrcid, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestGetOrcidDetailsEmptyToken( t *testing.T ) {
    expected := http.StatusBadRequest
    status, _ := client.GetOrcidDetails( cfg.Endpoint, goodOrcid, empty )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestGetOrcidDetailsBadToken( t *testing.T ) {
    expected := http.StatusForbidden
    status, _ := client.GetOrcidDetails( cfg.Endpoint, goodOrcid, badToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

//
// search ORCID tests
//

func TestSearchOrcidHappyDay( t *testing.T ) {

    expected := http.StatusOK
    status, orcids := client.SearchOrcid( cfg.Endpoint, goodSearch, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }

    ensureValidOrcidDetails( t, orcids )
}

func TestSearchOrcidEmptySearch( t *testing.T ) {
    expected := http.StatusBadRequest
    status, _ := client.SearchOrcid( cfg.Endpoint, empty, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestSearchOrcidNotFoundSearch( t *testing.T ) {
    expected := http.StatusNotFound
    status, _ := client.SearchOrcid( cfg.Endpoint, notFoundSearch, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestSearchOrcidEmptyToken( t *testing.T ) {
    expected := http.StatusBadRequest
    status, _ := client.SearchOrcid( cfg.Endpoint, goodSearch, empty )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestSearchOrcidBadToken( t *testing.T ) {
    expected := http.StatusForbidden
    status, _ := client.SearchOrcid( cfg.Endpoint, goodSearch, badToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

//
// set ORCID for user tests
//

func TestSetOrcidHappyDay( t *testing.T ) {
    expected := http.StatusOK
    status := client.SetOrcid( cfg.Endpoint, randomCid( ), randomOrcid( ), goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestSetOrcidDuplicate( t *testing.T ) {
    expected := http.StatusOK
    cid := randomCid( )
    orcid1 := randomOrcid( )
    orcid2 := randomOrcid( )

    status := client.SetOrcid( cfg.Endpoint, cid, orcid1, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }

    status = client.SetOrcid( cfg.Endpoint, cid, orcid2, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestSetOrcidEmptyId( t *testing.T ) {
    expected := http.StatusBadRequest
    status := client.SetOrcid( cfg.Endpoint, empty, goodOrcid, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestSetOrcidEmptyOrcid( t *testing.T ) {
    expected := http.StatusBadRequest
    status := client.SetOrcid( cfg.Endpoint, goodCid, empty, goodToken )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestSetOrcidEmptyToken( t *testing.T ) {
    expected := http.StatusBadRequest
    status := client.SetOrcid( cfg.Endpoint, goodCid, goodOrcid, empty )
    if status != expected {
        t.Fatalf( "Expected %v, got %v\n", expected, status )
    }
}

func TestSetOrcidBadToken( t *testing.T ) {
    expected := http.StatusForbidden
    status := client.SetOrcid( cfg.Endpoint, goodCid, goodOrcid, badToken )
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

func randomOrcid( ) string {

    // see the RNG
    rand.Seed( time.Now( ).UnixNano( ) )

    // list of possible characters
    possible := []rune( "0123456789" )
    return( fmt.Sprintf( "%s-%s-%s-%s", randomString( possible, 4 ), randomString( possible, 4 ),
        randomString( possible, 4 ), randomString( possible, 4 ) ) )
}

func randomCid( ) string {

    // see the RNG
    rand.Seed( time.Now( ).UnixNano( ) )

    // list of possible characters
    possible := []rune( "abcdefghijklmnopqrstuvwxyz0123456789" )

    return randomString( possible, 5 )
}

func randomString( possible []rune, sz int ) string {

    b := make( []rune, sz )
    for i := range b {
        b[ i ] = possible[ rand.Intn( len( possible ) ) ]
    }
    return string( b )
}

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

func ensureValidOrcidDetails( t *testing.T, orcids [] * api.OrcidDetails ) {
    for _, e := range orcids {
        //if emptyField( e.Id ) ||
        //        emptyField( e.Cid ) ||
        //        emptyField( e.Orcid ) ||
        //        emptyField( e.CreatedAt ) {
        if emptyField( e.Id ) ||
           emptyField( e.Uri ) ||
           //emptyField( e.DisplayName ) ||
           emptyField( e.FirstName ) ||
           emptyField( e.LastName ) {
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