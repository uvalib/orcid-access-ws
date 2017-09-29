package orcid

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"orcidaccessws/api"
	"orcidaccessws/logger"
    "orcidaccessws/config"
	"strings"
    "text/template"
    "bytes"
)

//
// log the contents of an activity update request
//
func logActivityUpdateRequest( activity api.ActivityUpdate ) {

	if config.Configuration.Debug {
       fmt.Println("UpdateCode:", activity.UpdateCode )
       fmt.Println("Work Title:", activity.Work.Title )
       fmt.Println("Work Abstract:", activity.Work.Abstract )
       fmt.Println("PublicationDate:", activity.Work.PublicationDate )
       fmt.Println("Url:", activity.Work.Url )
       fmt.Println("Authors:", activity.Work.Authors )
       fmt.Println("ResourceType:", activity.Work.ResourceType )
	}
}

func makeUpdateActivityBody( activity api.ActivityUpdate ) (string, error) {

   t, err := template.ParseFiles("data/work-activity-template.xml")
   if err != nil {
      logger.Log(fmt.Sprintf("ERROR: template parse error: %s", err))
      return "", err
   }

   // parse the publication date
   YYYY, MM, DD := splitDate( activity.Work.PublicationDate )

   // create our template data structure
   data := struct {
      Title             string
      Abstract          string
      ResourceType      string
      PublicationYear   string
      PublicationMonth  string
      PublicationDay    string
      Identifier        string
      Url               string
      Authors           []api.Person
   }{
      htmlEncodeString( activity.Work.Title ),
      htmlEncodeString( activity.Work.Abstract ),
      activity.Work.ResourceType,
      YYYY,
      MM,
      DD,
      "xxx",
      activity.Work.Url,
      htmlEncodePersonArray( api.SortPeople(activity.Work.Authors ) ),
   }

   var buffer bytes.Buffer
   err = t.Execute(&buffer, data)
   if err != nil {
      logger.Log(fmt.Sprintf("ERROR: template execute error: %s", err))
      return "", err
   }

   s := buffer.String()
   //s = filterUnacceptableCharacters(s)

   if config.Configuration.Debug {
      fmt.Printf("XML:\n%s\n", s)
   }
   return s, nil
}

func checkCommonResponse(body string) (int, error) {

	cr := orcidCommonResponse{}
	err := json.Unmarshal([]byte(body), &cr)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: json unmarshal: %s", err))
		return http.StatusInternalServerError, err
	}

	// check protocol version to ensure we know what to do with this
	if cr.Version != PROTOCOL_VERSION {
		logger.Log(fmt.Sprintf("ORCID protocol version not supported. Require: %s, received: %s", PROTOCOL_VERSION, cr.Version))
		return http.StatusHTTPVersionNotSupported, nil
	}

	// is there an error string
	if cr.Error.Value != "" {
		if strings.HasPrefix(cr.Error.Value, "Not found") == true {
			return http.StatusNotFound, nil
		}

		// not sure, just return a general error
		return http.StatusInternalServerError, errors.New(cr.Error.Value)
	}

	return http.StatusOK, nil
}

func transformDetailsResponse(profile *orcidProfile) *api.OrcidDetails {
	return constructDetails(profile)
}

func transformSearchResponse(search orcidResults) []*api.OrcidDetails {
	results := make([]*api.OrcidDetails, 0)
	for _, e := range search.Results {
		od := constructDetails(&e.Profile)
		od.Relevancy = fmt.Sprintf("%.6f", e.Relevancy.Value)
		results = append(results, od)
	}
	return (results)
}

func constructDetails(profile *orcidProfile) *api.OrcidDetails {

	od := new(api.OrcidDetails)

	od.Orcid = profile.Id.Id
	od.Uri = profile.Id.Uri
	od.DisplayName = profile.Bio.PersonalDetails.CreditName.Value
	od.FirstName = profile.Bio.PersonalDetails.GivenName.Value
	od.LastName = profile.Bio.PersonalDetails.FamilyName.Value
	od.Biography = profile.Bio.Biography.Value

	od.Keywords = make([]string, 0)
	for _, e := range profile.Bio.Keywords.Keywords {
		od.Keywords = append(od.Keywords, e.Value)
	}

	od.ResearchUrls = make([]string, 0)
	for _, e := range profile.Bio.Urls.Urls {
		od.ResearchUrls = append(od.ResearchUrls, e.Url.Value)
	}

	return (od)
}

//
// when including content embedded in XML, we should HTML encode it.
//
func htmlEncodePersonArray(array []api.Person) []api.Person {

   encoded := make([]api.Person, len(array), len(array))
   for ix, value := range array {

      p := api.Person{
         Index:       value.Index,
         FirstName:   htmlEncodeString(value.FirstName),
         LastName:    htmlEncodeString(value.LastName),
      }
      encoded[ix] = p
   }
   return encoded
}

func htmlEncodeString(value string) string {
   // HTML encoding
   encoded := html.EscapeString(value)

   // encode percent characters
   encoded = strings.Replace(encoded, "%", "%25", -1)
   return encoded
}

//
// Split a date in the form YYYY-MM-DD into its components
//
func splitDate(date string) (string, string, string) {
   tokens := strings.Split(date, "-")
   var YYYY, MM, DD string
   if len(tokens) > 0 {
      YYYY = tokens[0]
   }

   if len(tokens) > 1 {
      MM = tokens[1]
   }

   if len(tokens) > 2 {
      DD = tokens[2]
   }
   return YYYY, MM, DD
}

