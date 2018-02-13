package main

import (
	"fmt"
	"log"
	"net/http"
	"orcidaccessws/config"
	"orcidaccessws/dao"
	"orcidaccessws/handlers"
	"orcidaccessws/logger"
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

	// setup router and serve...
	router := NewRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Configuration.ServicePort), router))
}

//
// end of file
//
