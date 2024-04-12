package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/api"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/authtoken"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/config"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/dao"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/orcid"
	"io"
	"io/ioutil"
	"net/http"
)

var emptyUpdateCode = ""

// UpdateActivity -- update activity handler
func UpdateActivity(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	token := r.URL.Query().Get("auth")

	// parameters OK?
	if isEmpty(id) || isEmpty(token) {
		status := http.StatusBadRequest
		encodeUpdateActivityResponse(w, status, http.StatusText(status), emptyUpdateCode)
		return
	}

	// validate the token
	if authtoken.Validate(config.Configuration.SharedSecret, token) == false {
		status := http.StatusForbidden
		encodeUpdateActivityResponse(w, status, http.StatusText(status), emptyUpdateCode)
		return
	}

	decoder := json.NewDecoder(r.Body)
	activity := api.ActivityUpdate{}

	if err := decoder.Decode(&activity); err != nil {
		logger.Log(fmt.Sprintf("ERROR: decoding update activity request payload: %s", err))
		status := http.StatusBadRequest
		encodeUpdateActivityResponse(w, status,
			fmt.Sprintf("%s (%s)", http.StatusText(status), err.Error()),
			emptyUpdateCode)
		return
	}

	defer io.Copy(ioutil.Discard, r.Body)
	defer r.Body.Close()

	if err := validateRequestPayload(activity); err != nil {
		logger.Log(fmt.Sprintf("ERROR: invalid request payload: %s", err))
		status := http.StatusBadRequest
		encodeUpdateActivityResponse(w, status,
			fmt.Sprintf("%s (%s)", http.StatusText(status), err.Error()),
			emptyUpdateCode)
		return
	}

	// get the users ORCID attributes
	attributes, err := dao.Store.GetOrcidAttributesByCid(id)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: %s", err.Error()))
		status := http.StatusInternalServerError
		encodeUpdateActivityResponse(w, status,
			fmt.Sprintf("%s (%s)", http.StatusText(status), err.Error()),
			emptyUpdateCode)
		return
	}

	// we did not find the item, return 404
	if attributes == nil || len(attributes) == 0 {
		status := http.StatusNotFound
		encodeUpdateActivityResponse(w, status, http.StatusText(status), emptyUpdateCode)
		return
	}

	// verify the attributes are sufficient for our needs
	if err := validateOrcidAttributes(*attributes[0]); err != nil {
		logger.Log(fmt.Sprintf("ERROR: invalid ORCID attributes for cid %s: %s", id, err))
		status := http.StatusBadRequest
		encodeUpdateActivityResponse(w, status,
			fmt.Sprintf("%s (%s)", http.StatusText(status), err.Error()),
			emptyUpdateCode)
		return
	}

	// update the activity
	updateCode, status, err := orcid.UpdateOrcidActivity(attributes[0].Orcid, attributes[0].OauthAccessToken, activity)

	// the token might be expired, lets try to renew it
	if false { //status == http.StatusUnauthorized {
		var newAccessToken = ""
		var newRefreshToken = ""
		// renew the access token...
		newAccessToken, newRefreshToken, status, err = orcid.RenewAccessToken(attributes[0].OauthRefreshToken)
		if status == http.StatusOK {
			attributes[0].OauthAccessToken = newAccessToken
			attributes[0].OauthRefreshToken = newRefreshToken
			// save the new tokens
			err = dao.Store.SetOrcidAttributesByCid(id, *attributes[0])

			// if successful, retry the activity update
			if err == nil {
				updateCode, status, err = orcid.UpdateOrcidActivity(attributes[0].Orcid, attributes[0].OauthAccessToken, activity)
			} else {
				logger.Log(fmt.Sprintf("ERROR: %s", err.Error()))
				status = http.StatusInternalServerError
			}
		}
	}

	// we did got an error, return it
	if status != http.StatusOK {
		encodeUpdateActivityResponse(w, status,
			fmt.Sprintf("%s (%s)", http.StatusText(status), err), emptyUpdateCode)
		return
	}

	encodeUpdateActivityResponse(w, status, http.StatusText(status), updateCode)
}

// basic validation that the required fields for the activity update request exist
func validateRequestPayload(activity api.ActivityUpdate) error {

	if len(activity.Work.Title) == 0 {
		return errors.New("empty work title")
	}

	if len(activity.Work.ResourceType) == 0 {
		return errors.New("empty work resource type")
	}

	if len(activity.Work.URL) == 0 {
		return errors.New("empty work url")
	}

	return nil
}

// validation that the necessary ORCID attributes exist before we use them
func validateOrcidAttributes(attributes api.OrcidAttributes) error {

	if len(attributes.Orcid) == 0 {
		return errors.New("blank ORCID attribute")
	}

	if len(attributes.OauthAccessToken) == 0 {
		return errors.New("blank OAuth access token")
	}

	if len(attributes.OauthRefreshToken) == 0 {
		return errors.New("blank OAuth refresh token")
	}

	return nil
}

//
// end of file
//
