package logger

import (
    "log"
)

func Log( msg string ) {
    log.Printf( "ORCIDACCESS: %s", msg )
}