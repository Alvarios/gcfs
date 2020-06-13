package database

import (
	"github.com/Alvarios/gcfs/config"
	"github.com/couchbase/gocb/v2"
	"testing"
)

func TestConnect(t *testing.T) {
	// Allow testing with already configured instances.
	config.LoadConfigForTest(config.Configuration{})
	err := Connect()

	if err != nil {
		t.Errorf("cannot connect to database instance : %s", err.Error())
	}

	_, err = Bucket.Ping(
		&gocb.PingOptions{
			ReportID:     "medication",
			ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeKeyValue},
		},
	)
	if err != nil {
		t.Errorf("cannot ping database instance : %s", err.Error())
		return
	}

	err = Cluster.Close(nil)
	if err != nil {
		t.Errorf("cannot close database instance : %s", err.Error())
		return
	}
}