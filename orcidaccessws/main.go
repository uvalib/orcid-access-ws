package main

import (
	"fmt"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/config"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/dao"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/handlers"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
	"log"
	"net/http"
	"time"
)

func main() {

	logger.Log(fmt.Sprintf("===> version: '%s' <===", handlers.Version()))

	// create the storage singleton
	err := dao.NewDatastore( )
	if err != nil {
		log.Fatal(err)
	}

	// setup router and server...
	serviceTimeout := 15 * time.Second
	router := NewRouter()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Configuration.ServicePort),
		Handler:      router,
		ReadTimeout:  serviceTimeout,
		WriteTimeout: serviceTimeout,
	}
	log.Fatal(server.ListenAndServe())
}

//
// end of file
//
