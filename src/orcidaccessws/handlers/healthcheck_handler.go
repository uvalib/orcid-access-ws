package handlers

import (
    "net/http"
    //"orcidaccessws/api"
    "orcidaccessws/dao"
    //"orcidaccessws/orcid"
)

func HealthCheck( w http.ResponseWriter, r *http.Request ) {

    // update the statistics
    Statistics.RequestCount++
    Statistics.HeartbeatCount++

    status := http.StatusOK
    db_err := dao.Database.Check( )
    orcid_err := (error)( nil )//orcid.GetStatus( )

    var db_msg, orcid_msg string

    if db_err != nil || orcid_err != nil {

        status = http.StatusInternalServerError

        if db_err != nil {
            db_msg = db_err.Error( )
        }

        if orcid_err != nil {
            orcid_msg = orcid_err.Error( )
        }
    }

    encodeHealthCheckResponse( w, status, db_msg, orcid_msg )
}