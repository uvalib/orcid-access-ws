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

func encodeOrcidResponse(w http.ResponseWriter, status int, message string, orcids []*api.Orcid) {

	logger.Log(fmt.Sprintf("encodeOrcidResponse status: %d (%s)", status, message))
	jsonAttributes(w)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.OrcidResponse{Status: status, Message: message, Orcids: orcids}); err != nil {
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

func encodeHealthCheckResponse(w http.ResponseWriter, status int, dbmsg string, orcidmsg string) {

	db_healthy, orcid_healthy := true, true
	if len(dbmsg) != 0 {
		db_healthy = false
	}
	if len(orcidmsg) != 0 {
		orcid_healthy = false
	}
	jsonAttributes(w)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(api.HealthCheckResponse{
		DbCheck:    api.HealthCheckResult{Healthy: db_healthy, Message: dbmsg},
		OrcidCheck: api.HealthCheckResult{Healthy: orcid_healthy, Message: orcidmsg}}); err != nil {
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
	if err := json.NewEncoder(w).Encode(api.RuntimeResponse{Version: version, CpuCount: cpus, GoRoutineCount: goroutines, ObjectCount: heapcount, AllocatedMemory: alloc}); err != nil {
		log.Fatal(err)
	}
}

func jsonAttributes(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func nonEmpty(param string) bool {
	return len(strings.TrimSpace(param)) != 0
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
