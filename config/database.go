package config

import "github.com/couchbase/gocb/v2"

type Database struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	BucketName string `json:"bucket_name"`
	Bucket *gocb.Bucket `json:"bucket"`
}
