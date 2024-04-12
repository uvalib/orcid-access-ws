package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
)

// Config -- our configuration structure
type Config struct {
	ServicePort string

	// database attributes
	DbSecure     string // do we use TLS
	DbHost       string // hostname of database server
	DbName       string // database name
	DbUser       string // database user name
	DbPassphrase string // database user password
	DbTimeout    string // connection/read/write timeout

	// ORCID attributes
	OrcidPublicURL    string
	OrcidSecureURL    string
	OrcidOauthURL     string
	OrcidClientID     string
	OrcidClientSecret string

	// token authentication
	SharedSecret string

	// Service timeout
	ServiceTimeout int

	// diagnostic only
	Debug bool
}

// Configuration -- our configuration instance
var Configuration = loadConfig()

func loadConfig() Config {

	// default value for the database timeout
	c := Config{DbTimeout: "10s"}

	// process command line flags and setup configuration
	flag.StringVar(&c.ServicePort, "port", "8080", "The service listen port")
	flag.StringVar(&c.DbSecure, "dbsecure", "false", "Use TLS for the database connection")
	flag.StringVar(&c.DbHost, "dbhost", "mysqldev.lib.virginia.edu:3306", "The database server hostname:port")
	flag.StringVar(&c.DbName, "dbname", "orcidaccess_development", "The database name")
	flag.StringVar(&c.DbUser, "dbuser", "orcidaccess", "The database username")
	flag.StringVar(&c.DbPassphrase, "dbpassword", "", "The database passphrase")
	flag.StringVar(&c.OrcidPublicURL, "orcidpublicurl", "https://pub.orcid.org/v3.0", "The ORCID service public URL")
	flag.StringVar(&c.OrcidSecureURL, "orcidsecureurl", "https://api.sandbox.orcid.org/v3.0", "The ORCID service secure URL")
	flag.StringVar(&c.OrcidOauthURL, "orcidoauthurl", "https://sandbox.orcid.org", "The ORCID service OAuth URL")
	flag.StringVar(&c.OrcidClientID, "orcidclientid", "client-id", "The ORCID client identifier")
	flag.StringVar(&c.OrcidClientSecret, "orcidclientsecret", "client-secret", "The ORCID client secret")
	flag.IntVar(&c.ServiceTimeout, "timeout", 15, "The external service timeout in seconds")
	flag.StringVar(&c.SharedSecret, "secret", "", "The JWT shared secret")
	flag.BoolVar(&c.Debug, "debug", false, "Enable debugging")

	flag.Parse()

	logger.Log(fmt.Sprintf("INFO: ServicePort:         %s", c.ServicePort))
	logger.Log(fmt.Sprintf("INFO: DbSecure:            %s", c.DbSecure))
	logger.Log(fmt.Sprintf("INFO: DbHost:              %s", c.DbHost))
	logger.Log(fmt.Sprintf("INFO: DbName:              %s", c.DbName))
	logger.Log(fmt.Sprintf("INFO: DbUser:              %s", c.DbUser))
	logger.Log(fmt.Sprintf("INFO: DbPassphrase:        %s", strings.Repeat("*", len(c.DbPassphrase))))
	logger.Log(fmt.Sprintf("INFO: DbTimeout:           %s", c.DbTimeout))
	logger.Log(fmt.Sprintf("INFO: OrcidPublicURL:      %s", c.OrcidPublicURL))
	logger.Log(fmt.Sprintf("INFO: OrcidSecureURL:      %s", c.OrcidSecureURL))
	logger.Log(fmt.Sprintf("INFO: OrcidOauthURL:       %s", c.OrcidOauthURL))
	logger.Log(fmt.Sprintf("INFO: OrcidClientID:       %s", c.OrcidClientID))
	logger.Log(fmt.Sprintf("INFO: OrcidClientSecret:   %s", strings.Repeat("*", len(c.OrcidClientSecret))))
	logger.Log(fmt.Sprintf("INFO: SharedSecret:        %s", strings.Repeat("*", len(c.SharedSecret))))
	logger.Log(fmt.Sprintf("INFO: ServiceTimeout:      %d", c.ServiceTimeout))
	logger.Log(fmt.Sprintf("INFO: Debug:               %t", c.Debug))

	return c
}

//
// end of file
//
