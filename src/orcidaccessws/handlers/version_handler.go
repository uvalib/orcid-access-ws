package handlers

import (
	"net/http"
)

func VersionGet(w http.ResponseWriter, r *http.Request) {
	// create response
	encodeVersionResponse(w, http.StatusOK, Version())
}
