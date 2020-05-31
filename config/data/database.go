package data

type Database struct {
	Address  string `json:"address"  default:"127.0.0.1"`
	Username string `json:"username"  default:""`
	Password string `json:"password"  default:""`
	BucketName string `json:"bucket_name" default:"metadata"`
}
