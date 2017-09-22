package handlers

import (
	"net/http"
	//  "github.com/gorilla/mux"
	"fmt"
	"orcidaccessws/authtoken"
	"orcidaccessws/config"
	"orcidaccessws/logger"
	"orcidaccessws/orcid"
)

const DEFAULT_SEARCH_START_IX = "0"
const DEFAULT_SEARCH_MAX_RESULTS = "50"

func SearchOrcid(w http.ResponseWriter, r *http.Request) {

	//vars := mux.Vars( r )
	query := r.URL.Query().Get("q")
	token := r.URL.Query().Get("auth")
	start := r.URL.Query().Get("start")
	count := r.URL.Query().Get("max")

	// update the statistics
	Statistics.RequestCount++
	Statistics.SearchOrcidDetailsCount++

	// parameters OK?
	if isEmpty(query) || isEmpty(token) {
		status := http.StatusBadRequest
		encodeOrcidSearchResponse(w, status, http.StatusText(status), nil, 0, 0, 0)
		return
	}

	// check the supplied parameters and set defaults as necessary
	if isEmpty(start) {
		start = DEFAULT_SEARCH_START_IX
	}
	if isEmpty(count) {
		count = DEFAULT_SEARCH_MAX_RESULTS
	}

	// validate parameters as necessary
	if isNumeric(start) == false || isNumeric(count) == false {
		status := http.StatusBadRequest
		encodeOrcidSearchResponse(w, status, http.StatusText(status), nil, 0, 0, 0)
		return
	}

	// validate the token
	if authtoken.Validate(config.Configuration.AuthTokenEndpoint, "getorcid", token, config.Configuration.Timeout) == false {
		status := http.StatusForbidden
		encodeOrcidSearchResponse(w, status, http.StatusText(status), nil, 0, 0, 0)
		return
	}

	// get the ORCID details
	orcids, total, status, err := orcid.SearchOrcid(query, start, count)

	// we got an error, return it
	if err != nil {
		encodeOrcidSearchResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("%s (%s)", http.StatusText(http.StatusInternalServerError), err), nil, 0, 0, 0)
		return
	}

	logger.Log(fmt.Sprintf("ORCID search: %d result(s) located", len(orcids)))

	// everything OK but found no items
	if len(orcids) == 0 {
		status = http.StatusNotFound
	} else {
		status = http.StatusOK
	}

	encodeOrcidSearchResponse(w, status, http.StatusText(status), orcids, asNumeric(start), len(orcids), total)
}
