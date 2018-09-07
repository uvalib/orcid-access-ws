package config

import (
	"flag"
	"fmt"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
	"strings"
)

//
// Config -- our configuration structure
//
type Config struct {
	ServicePort string

	// database attributes
	DbHost       string
	DbName       string
	DbUser       string
	DbPassphrase string
	DbTimeout    string

	// ORCID attributes
	OrcidPublicURL    string
	OrcidSecureURL    string
	OrcidOauthURL     string
	OrcidClientID     string
	OrcidClientSecret string

	// token authentication
	AuthTokenEndpoint string
	ServiceTimeout    int

	// diagnostic only
	Debug bool
}

//
// Configuration -- our configuration instance
//
var Configuration = loadConfig()

func loadConfig() Config {

	// default value for the database timeout
	c := Config{ DbTimeout: "15s" }

	// process command line flags and setup configuration
	flag.StringVar(&c.ServicePort, "port", "8080", "The service listen port")
	flag.StringVar(&c.DbHost, "dbhost", "mysqldev.lib.virginia.edu:3306", "The database server hostname:port")
	flag.StringVar(&c.DbName, "dbname", "orcidaccess_development", "The database name")
	flag.StringVar(&c.DbUser, "dbuser", "orcidaccess", "The database username")
	flag.StringVar(&c.DbPassphrase, "dbpassword", "", "The database passphrase")
	flag.StringVar(&c.OrcidPublicURL, "orcidpublicurl", "https://pub.orcid.org/v1.2", "The ORCID service public URL")
	flag.StringVar(&c.OrcidSecureURL, "orcidsecureurl", "https://api.sandbox.orcid.org/v2.0", "The ORCID service secure URL")
	flag.StringVar(&c.OrcidOauthURL, "orcidoauthurl", "https://sandbox.orcid.org", "The ORCID service OAuth URL")
	flag.StringVar(&c.OrcidClientID, "orcidclientid", "client-id", "The ORCID client identifier")
	flag.StringVar(&c.OrcidClientSecret, "orcidclientsecret", "client-secret", "The ORCID client secret")
	flag.IntVar(&c.ServiceTimeout, "timeout", 15, "The external service timeout in seconds")
	flag.StringVar(&c.AuthTokenEndpoint, "tokenauth", "http://docker1.lib.virginia.edu:8200", "The token authentication endpoint")
	flag.BoolVar(&c.Debug, "debug", false, "Enable debugging")

	flag.Parse()

	logger.Log(fmt.Sprintf("ServicePort:         %s", c.ServicePort))
	logger.Log(fmt.Sprintf("DbHost:              %s", c.DbHost))
	logger.Log(fmt.Sprintf("DbName:              %s", c.DbName))
	logger.Log(fmt.Sprintf("DbUser:              %s", c.DbUser))
	logger.Log(fmt.Sprintf("DbPassphrase:        %s", strings.Repeat("*", len(c.DbPassphrase))))
	logger.Log(fmt.Sprintf("DbTimeout:           %s", c.DbTimeout))
	logger.Log(fmt.Sprintf("OrcidPublicURL:      %s", c.OrcidPublicURL))
	logger.Log(fmt.Sprintf("OrcidSecureURL:      %s", c.OrcidSecureURL))
	logger.Log(fmt.Sprintf("OrcidOauthURL:       %s", c.OrcidOauthURL))
	logger.Log(fmt.Sprintf("OrcidClientID:       %s", c.OrcidClientID))
	logger.Log(fmt.Sprintf("OrcidClientSecret:   %s", strings.Repeat("*", len(c.OrcidClientSecret))))
	logger.Log(fmt.Sprintf("AuthTokenEndpoint    %s", c.AuthTokenEndpoint))
	logger.Log(fmt.Sprintf("ServiceTimeout:      %d", c.ServiceTimeout))
	logger.Log(fmt.Sprintf("Debug                %t", c.Debug))

	return c
}

//
// end of file
//