package main

import (
    "fmt"
    "log"
    "net/http"
    "orcidaccessws/api"
    "orcidaccessws/config"
    "orcidaccessws/logger"
)

var statistics = api.Statistics{ }

func main( ) {

    logger.Log( fmt.Sprintf( "===> version: '%s' <===", Version( ) ) )

	// setup router and serve...
    router := NewRouter( )
    log.Fatal( http.ListenAndServe( fmt.Sprintf( ":%s", config.Configuration.ServicePort ), router ) )
}

