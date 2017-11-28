package handlers

import (
	"log"
	//"fmt"
	"encoding/json"
	"net/http"
	"orcidaccessws/api"
	"strings"
	//"orcidaccessws/mapper"
	//"orcidaccessws/logger"
	"fmt"
	"orcidaccessws/logger"
	"strconv"
)

func encodeStandardResponse(w http.ResponseWriter, status int, message string) {

	logger.Log(fmt.Sprintf("encodeStandardResponse status: %d (%s)", status, message))
	jsonAttributes(w)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.StandardResponse{Status: status, Message: message}); err != nil {
		log.Fatal(err)
	}
}

func encodeOrcidAttributesResponse(w http.ResponseWriter, status int, message string, attributes []*api.OrcidAttributes) {

	logger.Log(fmt.Sprintf("encodeOrcidAttributesResponse status: %d (%s)", status, message))
	jsonAttributes(w)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.OrcidAttributesResponse{Status: status, Message: message, Attributes: attributes}); err != nil {
		log.Fatal(err)
	}
}

func encodeUpdateActivityResponse(w http.ResponseWriter, status int, message string, updateCode string) {

	logger.Log(fmt.Sprintf("encodeUpdateActivityResponse status: %d (%s)", status, message))
	jsonAttributes(w)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.UpdateActivityResponse{Status: status, Message: message, UpdateCode: updateCode}); err != nil {
		log.Fatal(err)
	}
}

func encodeOrcidDetailsResponse(w http.ResponseWriter, status int, message string, details *api.OrcidDetails) {

	logger.Log(fmt.Sprintf("encodeOrcidDetailsResponse status: %d (%s)", status, message))
	jsonAttributes(w)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.OrcidDetailsResponse{Status: status, Message: message, Details: details}); err != nil {
		log.Fatal(err)
	}
}

func encodeOrcidSearchResponse(w http.ResponseWriter, status int, message string, results []*api.OrcidDetails,
	start int, count int, total int) {

	logger.Log(fmt.Sprintf("encodeOrcidSearchResponse status: %d (%s)", status, message))
	jsonAttributes(w)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.OrcidSearchResponse{Status: status, Message: message, Results: results,
		Start: start, Count: count, Total: total}); err != nil {
		log.Fatal(err)
	}
}

func encodeHealthCheckResponse(w http.ResponseWriter, status int, dbMsg string, orcidPublicMsg string, orcidSecureMsg string) {

	dbHealthy, orcidPublicHealthy, orcidSecureHealthy := true, true, true
	if len(dbMsg) != 0 {
		dbHealthy = false
	}
	if len(orcidPublicMsg) != 0 {
		orcidPublicHealthy = false
	}
	if len(orcidSecureMsg) != 0 {
		orcidSecureHealthy = false
	}

	jsonAttributes(w)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.HealthCheckResponse{
		DbCheck:          api.HealthCheckResult{Healthy: dbHealthy, Message: dbMsg},
		OrcidPublicCheck: api.HealthCheckResult{Healthy: orcidPublicHealthy, Message: orcidPublicMsg},
		OrcidSecureCheck: api.HealthCheckResult{Healthy: orcidSecureHealthy, Message: orcidSecureMsg}}); err != nil {
		log.Fatal(err)
	}
}

func encodeStatsResponse(w http.ResponseWriter, statistics api.Statistics) {

	status := http.StatusOK

	jsonAttributes(w)
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(api.StatisticsResponse{Status: status, Message: http.StatusText(status), Details: statistics}); err != nil {
		log.Fatal(err)
	}

}

func encodeVersionResponse(w http.ResponseWriter, status int, version string) {
	jsonAttributes(w)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.VersionResponse{Version: version}); err != nil {
		log.Fatal(err)
	}
}

func encodeRuntimeResponse(w http.ResponseWriter, status int, version string, cpus int, goroutines int, heapcount uint64, alloc uint64) {
	jsonAttributes(w)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.RuntimeResponse{Version: version, CPUCount: cpus, GoRoutineCount: goroutines, ObjectCount: heapcount, AllocatedMemory: alloc}); err != nil {
		log.Fatal(err)
	}
}

func jsonAttributes(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func isEmpty(param string) bool {
	return len(strings.TrimSpace(param)) == 0
}

func isNumeric(param string) bool {
	if _, err := strconv.Atoi(param); err == nil {
		return true
	}
	return false
}

func asNumeric(param string) int {
	res, _ := strconv.Atoi(param)
	return (res)
}

//
// end of file
//
