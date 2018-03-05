package main

import (
	"fmt"
	"log"
	"net/http"
	"orcidaccessws/config"
	"orcidaccessws/dao"
	"orcidaccessws/handlers"
	"orcidaccessws/logger"
	"time"
)

func main() {

	logger.Log(fmt.Sprintf("===> version: '%s' <===", handlers.Version()))

	// access the database
	connectStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?allowOldPasswords=1&sql_notes=false", config.Configuration.DbUser,
		config.Configuration.DbPassphrase, config.Configuration.DbHost, config.Configuration.DbName)

	err := dao.NewDB(connectStr)
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
