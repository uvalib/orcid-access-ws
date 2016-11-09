package config

import (
    "flag"
    "fmt"
    "orcidaccessws/logger"
)

type Config struct {
    ServicePort         string
    DbHost              string
    DbName              string
    DbUser              string
    DbPassphrase        string
    OrcidServiceUrl     string
    OrcidServiceTimeout int
    OrcidUser           string
    OrcidPassphrase     string
    AuthTokenEndpoint   string
    Debug               bool
}

var Configuration = LoadConfig( )

func LoadConfig( ) Config {

    c := Config{ }

    // process command line flags and setup configuration
    flag.StringVar( &c.ServicePort, "port", "8080", "The service listen port" )
    flag.StringVar( &c.DbHost, "dbhost", "mysqldev.lib.virginia.edu:3306", "The database server hostname:port" )
    flag.StringVar( &c.DbName, "dbname", "orcidaccess_development", "The database name" )
    flag.StringVar( &c.DbUser, "dbuser", "orcidaccess", "The database username" )
    flag.StringVar( &c.DbPassphrase, "dbpassword", "dbpassword", "The database passphrase")
    flag.StringVar( &c.OrcidServiceUrl, "orcidurl", "https://ezid.cdlib.org", "The ORCID service URL" )
    flag.IntVar( &c.OrcidServiceTimeout, "timeout", 10, "The service timeout (in seconds)")
    flag.StringVar( &c.OrcidUser, "orciduser", "apitest", "The EZID service username" )
    flag.StringVar( &c.OrcidPassphrase, "orcidpassword", "apitest", "The ORCID service passphrase")
    flag.StringVar( &c.AuthTokenEndpoint, "tokenauth", "http://docker1.lib.virginia.edu:8200", "The token authentication endpoint")
    flag.BoolVar( &c.Debug, "debug", false, "Enable debugging")

    flag.Parse( )

    logger.Log( fmt.Sprintf( "ServicePort:         %s", c.ServicePort ) )
    logger.Log( fmt.Sprintf( "DbHost:              %s", c.DbHost ) )
    logger.Log( fmt.Sprintf( "DbName:              %s", c.DbName ) )
    logger.Log( fmt.Sprintf( "DbUser:              %s", c.DbUser ) )
    logger.Log( fmt.Sprintf( "DbPassphrase:        %s", c.DbPassphrase ) )
    logger.Log( fmt.Sprintf( "OrcidServiceUrl:     %s", c.OrcidServiceUrl ) )
    logger.Log( fmt.Sprintf( "OrcidServiceTimeout: %d", c.OrcidServiceTimeout ) )
    logger.Log( fmt.Sprintf( "OrcidUser:           %s", c.OrcidUser ) )
    logger.Log( fmt.Sprintf( "OrcidPassphrase:     %s", c.OrcidPassphrase ) )
    logger.Log( fmt.Sprintf( "AuthTokenEndpoint    %s", c.AuthTokenEndpoint ) )
    logger.Log( fmt.Sprintf( "Debug                %t", c.Debug ) )

    return c
}
