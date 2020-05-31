package api

import (
	"encoding/json"
	"gcfs"
	"gcfs/api/responses"
	"github.com/couchbase/gocb/v2"
	"log"
	"net/http"
)

func PingCouchbase(w http.ResponseWriter, _ *http.Request) {
	// write header first to avoid flush error.
	w.WriteHeader(http.StatusOK)

	pings, err := gcfs.Cluster.Bucket("metadata").Ping(
		&gocb.PingOptions{
			ReportID:     "medication",
			ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeKeyValue},
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for service, pingReports := range pings.Services {
		if service != gocb.ServiceTypeKeyValue {
			log.Printf("unwanted service type !")
			w.WriteHeader(http.StatusInternalServerError)
			break
		}

		for _, pingReport := range pingReports {
			if pingReport.State != gocb.PingStateOk {
				log.Printf(
					"Node %s at remote %s is not OK, error: %s, latency: %s\n",
					pingReport.ID, pingReport.Remote, pingReport.Error, pingReport.Latency.String(),
				)
			} else {
				log.Printf(
					"Node %s at remote %s is OK, latency: %s\n",
					pingReport.ID, pingReport.Remote, pingReport.Latency.String(),
				)
			}
		}
	}

	b, err := json.Marshal(pings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	_, err = w.Write(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
