package handlers

import (
	"net/http"
	//  "github.com/gorilla/mux"
	//"fmt"
	"orcidaccessws/authtoken"
	"orcidaccessws/config"
	//"orcidaccessws/logger"
	//"orcidaccessws/orcid"
)

const defaulsSearchStartIx = "0"
const defaultSearchMaxResults = "50"

//
// SearchOrcid -- the search orcid handler
//
func SearchOrcid(w http.ResponseWriter, r *http.Request) {

	//vars := mux.Vars( r )
	query := r.URL.Query().Get("q")
	token := r.URL.Query().Get("auth")
	start := r.URL.Query().Get("start")
	count := r.URL.Query().Get("max")

	// parameters OK?
	if isEmpty(query) || isEmpty(token) {
		status := http.StatusBadRequest
		encodeOrcidSearchResponse(w, status, http.StatusText(status), nil, 0, 0, 0)
		return
	}

	// check the supplied parameters and set defaults as necessary
	if isEmpty(start) {
		start = defaulsSearchStartIx
	}
	if isEmpty(count) {
		count = defaultSearchMaxResults
	}

	// validate parameters as necessary
	if isNumeric(start) == false || isNumeric(count) == false {
		status := http.StatusBadRequest
		encodeOrcidSearchResponse(w, status, http.StatusText(status), nil, 0, 0, 0)
		return
	}

	// validate the token
	if authtoken.Validate(config.Configuration.AuthTokenEndpoint, "getorcid", token, config.Configuration.ServiceTimeout) == false {
		status := http.StatusForbidden
		encodeOrcidSearchResponse(w, status, http.StatusText(status), nil, 0, 0, 0)
		return
	}

	//
	// not implemented as we have moved to the 2.0 API which supports different behavior
	//
	status := http.StatusNotImplemented
	encodeOrcidSearchResponse(w, status, http.StatusText(status), nil, 0, 0, 0)

	// get the ORCID details
	//orcids, total, status, err := orcid.SearchOrcid(query, start, count)

	// we got an error, return it
	//if err != nil {
	//   encodeOrcidSearchResponse(w, status,
	//      fmt.Sprintf("%s (%s)", http.StatusText(status), err), nil, 0, 0, 0)
	//   return
	//}

	//logger.Log(fmt.Sprintf("ORCID search: %d result(s) located", len(orcids)))

	// everything OK but found no items
	//if len(orcids) == 0 {
	//   status = http.StatusNotFound
	//} else {
	//   status = http.StatusOK
	//}

	//encodeOrcidSearchResponse(w, status, http.StatusText(status), orcids, asNumeric(start), len(orcids), total)
}

//
// end of file
//
