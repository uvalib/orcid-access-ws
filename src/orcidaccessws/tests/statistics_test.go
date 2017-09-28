package test

import (
   "net/http"
   "orcidaccessws/client"
   "testing"
)

//
// statistics tests
//

func TestStatistics(t *testing.T) {
   expected := http.StatusOK
   status, stats := client.Statistics(cfg.Endpoint)
   if status != expected {
      t.Fatalf("Expected %v, got %v\n", expected, status)
   }

   if stats.RequestCount == 0 {
      t.Fatalf("Expected non-zero request count\n")
   }
}

//
// end of file
//